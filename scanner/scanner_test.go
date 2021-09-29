package scanner

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	ast := assert.New(t)

	r := &MockReader{
		mask:   2,
		fields: "country,isp",
		data:   "中国:电信",
	}
	b := make([]byte, 0)
	buf := bytes.NewBuffer(b)

	scanner := NewScanner(r, buf)
	ast.NotNil(scanner)

	scanner.ScanAndMerge(false)
	scanResult := string(buf.Bytes())
	ast.True(strings.Contains(scanResult, "# IPVersion: 1"))
	ast.True(strings.Contains(scanResult, "# Fields: country,isp"))
	// merge all
	ast.True(strings.Contains(scanResult, "0.0.0.0/0\t中国:电信"))

	r.diffData = true
	buf.Reset()
	scanner.ScanAndMerge(false)
	scanResult = string(buf.Bytes())

	ast.True(strings.Contains(scanResult, "0.0.0.0/2\t0.0.0.0"))
	ast.True(strings.Contains(scanResult, "64.0.0.0/2\t64.0.0.0"))
	ast.True(strings.Contains(scanResult, "128.0.0.0/2\t128.0.0.0"))
	ast.True(strings.Contains(scanResult, "192.0.0.0/2\t192.0.0.0"))

	buf.Reset()
	scanner.ScanAndMerge(true)
	scanResult = string(buf.Bytes())
	ast.True(strings.Contains(scanResult, "# IPVersion: 2"))
	ast.True(strings.Contains(scanResult, "::/2\t::"))
	ast.True(strings.Contains(scanResult, "4000::/2\t4000::"))
	ast.True(strings.Contains(scanResult, "8000::/2\t8000::"))
	ast.True(strings.Contains(scanResult, "c000::/2\tc000::"))
}

type MockReader struct {
	mask     int
	fields   string
	data     string
	diffData bool
}

func (r *MockReader) LookupNetwork(ip net.IP) (*net.IPNet, string, error) {
	_, ipNet, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip.String(), r.mask))
	if r.diffData {
		return ipNet, ip.String(), nil
	}
	return ipNet, r.data, nil
}

func (r *MockReader) FieldsCollection() string {
	return r.fields
}
