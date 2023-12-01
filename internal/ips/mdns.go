/*
 * Copyright (c) 2023 shenjunzheng@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ips

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/internal/data"
	"github.com/sjzar/ips/internal/util"
)

// MDNSResolve resolves the given domain name using MDNS and formats the results in a table.
func (m *Manager) MDNSResolve(domain string) (string, error) {

	if m.mdns == nil {
		var err error
		m.mdns, err = NewMDNS(m.Conf)
		if err != nil {
			log.Errorf("Failed to initialize MDNS: %v", err)
			return "", err
		}
	}

	ret, err := m.mdns.Resolve(domain)
	if err != nil {
		log.Errorf("Failed to resolve domain: %v", err)
		return "", err
	}

	return m.MDNSFormatTable(ret), nil
}

// MDNSFormatTable formats the MDNS response into a readable table string.
func (m *Manager) MDNSFormatTable(ret *MdnsResponse) string {
	cnameMap := map[string]struct{}{}
	ipMap := map[string]struct{}{}

	writer := &strings.Builder{}
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"GeoISP", "CNAME", "IP"})
	table.SetRowLine(true)
	table.SetAutoWrapText(false)

	for _, item := range ret.Items {
		geoISP := fmt.Sprintf("%s\n[%s]", item.IP,
			strings.Join(util.DeleteEmptyValue(strings.Split(item.GeoISP, ",")), m.Conf.TextValuesSep))
		for i := range item.CNAME {
			if info, err := m.parseDomain(item.CNAME[i][:len(item.CNAME[i])-1]); err == nil {
				if text, err := m.serializeDomainInfoToText(info); err == nil {
					item.CNAME[i] = text
				}
			}
		}
		if len(item.CNAME) > 0 {
			cnameMap[item.CNAME[len(item.CNAME)-1]] = struct{}{}
		}
		for i := range item.Result {
			ipMap[item.Result[i]] = struct{}{}
			if info, err := m.parseIP(item.Result[i]); err == nil {
				if text, err := m.serializeIPInfoToText(info); err == nil {
					item.Result[i] = text
				}
			}
		}
		table.Append([]string{
			geoISP,
			strings.Join(item.CNAME, "\n"),
			strings.Join(item.Result, "\n"),
		})
	}

	table.SetFooter([]string{"Total", fmt.Sprintf("%d", len(cnameMap)), fmt.Sprintf("%d", len(ipMap))})
	table.Render()

	return writer.String()
}

// MDNS represents a structure to hold MDNS client and configuration.
type MDNS struct {
	config *Config         // Configuration for MDNS.
	client *dns.Client     // DNS client for MDNS queries.
	views  map[string]View // Map of views for different geolocations.
}

// View represents a geolocation view with associated IP and GeoISP information.
type View struct {
	GeoIPS string // Geolocation ISP information.
	IP     net.IP // Associated IP address.
}

// MdnsResponse represents the response from an MDNS query.
type MdnsResponse struct {
	Domain string     // Queried domain.
	Items  []MdnsItem // List of MDNS items in the response.
}

// MdnsItem represents a single item in the MDNS response, containing relevant DNS information.
type MdnsItem struct {
	GeoISP string   // Geolocation ISP information.
	IP     string   // IP address.
	CNAME  []string // CNAME records.
	Result []string // IP resolution results.
}

// NewMDNS creates a new MDNS instance with the provided configuration.
func NewMDNS(conf *Config) (*MDNS, error) {
	client := new(dns.Client)
	client.Net = conf.DNSClientNet
	client.Timeout = time.Millisecond * time.Duration(conf.DNSClientTimeoutMs)
	client.SingleInflight = conf.DNSClientSingleInflight

	views, err := loadViews()
	if err != nil {
		log.Errorf("Failed to load views: %v", err)
		return nil, err
	}

	return &MDNS{
		config: conf,
		client: client,
		views:  views,
	}, nil
}

// Resolve performs the domain name resolution using the MDNS setup.
func (m *MDNS) Resolve(domain string) (*MdnsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.config.MDNSTimeoutS))
	defer cancel()

	retChan := make(chan MdnsItem, len(m.views))
	errChan := make(chan error, len(m.views))
	wg := sync.WaitGroup{}

	wg.Add(len(m.views))
	for _, view := range m.views {
		go func(ctx context.Context, view View) {
			defer wg.Done()
			cname, ips, err := m.do(ctx, m.config.MDNSExchangeAddress+":53", domain, view.IP)
			if err != nil {
				errChan <- err
				return
			}
			retChan <- MdnsItem{
				GeoISP: view.GeoIPS,
				IP:     view.IP.String(),
				CNAME:  cname,
				Result: ips,
			}
		}(ctx, view)
	}

	wg.Wait()
	close(retChan)
	close(errChan)

	if len(errChan) > 0 {
		var errs []string
		for e := range errChan {
			errs = append(errs, e.Error())
		}
		return nil, fmt.Errorf("multiple resolve errors: %s", strings.Join(errs, "; "))
	}

	ret := &MdnsResponse{
		Domain: domain,
		Items:  make([]MdnsItem, 0, len(m.views)),
	}

	for item := range retChan {
		ret.Items = append(ret.Items, item)
	}

	return ret, nil
}

// do performs the actual DNS query for the given domain and view.
func (m *MDNS) do(ctx context.Context, address, domain string, ip net.IP) (cname, ips []string, err error) {
	var msg *dns.Msg
	for i := 0; i <= m.config.MDNSRetryTimes; i++ {
		req := makeEDNSQueryMsg(domain, ip)
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			msg, _, err = m.client.ExchangeContext(ctx, req, address)
			if err == nil {
				cname, ips := parseMsg(msg)
				return cname, ips, nil
			}
		}
	}
	return nil, nil, err
}

// parseMsg parses the DNS message and extracts CNAME and IP information.
func parseMsg(msg *dns.Msg) ([]string, []string) {
	cname := make([]string, 0, len(msg.Answer))
	ips := make([]string, 0, len(msg.Answer))
	for _, ans := range msg.Answer {
		if record, ok := ans.(*dns.CNAME); ok {
			cname = append(cname, record.Target)
		}
		if record, ok := ans.(*dns.A); ok {
			ips = append(ips, record.A.String())
		}

		if record, ok := ans.(*dns.AAAA); ok {
			ips = append(ips, record.AAAA.String())
		}
	}
	return cname, ips
}

// makeEDNSQueryMsg constructs an EDNS query message for the given domain and IP.
func makeEDNSQueryMsg(domain string, ip net.IP) *dns.Msg {
	dnsType := dns.TypeA
	family := uint16(1)
	sourceNetmask := uint8(24)
	if netip.MustParseAddr(ip.String()).Is6() {
		dnsType = dns.TypeAAAA
		family = uint16(2)
		sourceNetmask = uint8(56)
	}

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dnsType)
	msg.RecursionDesired = true

	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_SUBNET)
	e.Code = dns.EDNS0SUBNET
	e.Family = family
	e.SourceNetmask = sourceNetmask
	e.SourceScope = 0
	e.Address = ip
	o.Option = append(o.Option, e)
	msg.Extra = append(msg.Extra, o)

	return msg
}

// loadViews loads geolocation views from a data source.
func loadViews() (map[string]View, error) {

	ret := make(map[string]View)
	scanner := bufio.NewScanner(strings.NewReader(data.MdnsView))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		split := strings.SplitN(string(line), "\t", 2)
		if len(split) < 2 {
			continue
		}
		ip, _, err := net.ParseCIDR(split[1])
		if err != nil {
			continue
		}
		ret[split[0]] = View{
			GeoIPS: split[0],
			IP:     ip,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}
