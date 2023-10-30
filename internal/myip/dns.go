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

package myip

import (
	"context"
	"net"
	"time"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/pkg/errors"
)

// DNSDetector is a struct for configuring DNS queries.
// It contains the domain to query, the DNS server to use, and the type of DNS query.
type DNSDetector struct {
	Domain    string
	Server    string
	QueryType string
}

// Discover performs a DNS query based on the DNSDetector configuration and returns the discovered IP.
// It supports A and TXT record queries.
func (d *DNSDetector) Discover(ctx context.Context, localAddr string) (net.IP, error) {
	var dialer *net.Dialer
	if len(localAddr) > 0 {
		addr, err := net.ResolveUDPAddr("udp", localAddr+":0")
		if err != nil {
			log.Debugf("Failed to resolve UDP address: %s", err)
			return nil, err
		}
		dialer = &net.Dialer{LocalAddr: addr, Timeout: time.Second * 10}
	}
	c := dns.Client{Dialer: dialer}
	m := dns.Msg{}
	dnsType := d.getQueryType()

	m.SetQuestion(d.Domain, dnsType)
	result, _, err := c.ExchangeContext(ctx, &m, d.Server)
	if err != nil {
		log.Debugf("dns query %s %s %s failed: %s", d.Domain, d.Server, d.QueryType, err)
		return nil, err
	}

	if len(result.Answer) == 0 {
		log.Debugf("no answer found in DNS response: %s", result.String())
		return nil, errors.ErrDiscoveryFailed
	}

	return d.parseAnswer(result.Answer[0])
}

// getQueryType returns the DNS query type based on the DNSDetector configuration.
func (d *DNSDetector) getQueryType() uint16 {
	if d.QueryType == "TXT" {
		return dns.TypeTXT
	}
	return dns.TypeA
}

// parseAnswer parses the DNS answer and returns the discovered IP.
func (d *DNSDetector) parseAnswer(answer dns.RR) (net.IP, error) {
	switch v := answer.(type) {
	case *dns.A:
		return v.A, nil
	case *dns.TXT:
		if len(v.Txt) == 0 {
			log.Debugf("no TXT record found in DNS answer: %s", v.String())
			return nil, errors.ErrDiscoveryFailed
		}
		return net.ParseIP(v.Txt[0]), nil
	default:
		log.Debugf("failed to parse DNS answer: %s", v.String())
		return nil, errors.ErrDiscoveryFailed
	}
}
