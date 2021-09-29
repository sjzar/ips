package scanner

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/sjzar/ips/iprange"
	"github.com/sjzar/ips/model"
	"github.com/sjzar/ips/packer"
	"github.com/sjzar/ips/reader"
)

// Scanner 扫描工具
// 扫描IP库中所有的CIDR分组信息，并输出到writer
type Scanner struct {
	reader reader.Reader
	writer io.Writer
}

// NewScanner 初始化扫描工具
func NewScanner(r reader.Reader, w io.Writer) *Scanner {
	return &Scanner{
		reader: r,
		writer: w,
	}
}

func (s *Scanner) ScanAndMerge(ipv6 bool) {

	bits, ipVersion := net.IPv4len, packer.IPv4
	if ipv6 {
		bits, ipVersion = net.IPv6len, packer.IPv6
	}

	s.MetaData(model.MetaScanTime, strconv.FormatInt(time.Now().Unix(), 10))
	s.MetaData(model.MetaIPVersion, strconv.Itoa(int(ipVersion)))
	s.MetaData(model.MetaFields, s.reader.FieldsCollection())

	var ipRange *iprange.IPRange
	dataTmp := "INIT_DATA"
	ipVersionCheck := false
	ZeroIP, markerIP := make(net.IP, bits), make(net.IP, bits)
	for {
		ipNet, data, err := s.reader.LookupNetwork(markerIP)
		if err != nil {
			log.Print(markerIP, err)
			return
		}

		if !ipVersionCheck {
			if (ipNet.IP.To4() == nil) != ipv6 {
				log.Fatal("db and scan version not match")
			}
			ipVersionCheck = true
		}

		if ipRange != nil {
			if data == dataTmp {
				ipRange.Join(ipNet)
				goto next
			} else {
				s.Print(ipRange, dataTmp)
			}
		}
		ipRange = iprange.NewIPRange(ipNet)
		ipRange.Start = markerIP
		dataTmp = data
	next:
		markerIP = iprange.NextIP(iprange.LastIP(ipNet))

		// loop end
		if markerIP.Equal(ZeroIP) {
			s.Print(ipRange, dataTmp)
			break
		}
	}
}

func (s *Scanner) MetaData(key, value string) {
	_, _ = fmt.Fprintf(s.writer, "# %s: %s\n", key, value)
}

func (s *Scanner) Print(ipRange *iprange.IPRange, data string) {
	for _, ipNet := range ipRange.IPNets() {
		_, _ = fmt.Fprintf(s.writer, "%s\t%s\n", ipNet, data)
	}
}
