package reader

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjzar/ips/model"
)

func TestReader(t *testing.T) {
	ast := assert.New(t)

	_, ipNet, err := net.ParseCIDR("0.0.0.0/32")
	ast.Nil(err)

	database := &MockDB{
		ipNet: ipNet,
		data: map[string]string{
			model.Country:  "country1",
			model.Province: "province1",
			model.City:     "city1",
			model.ISP:      "isp1",
		},
	}
	filter := NewKeyFieldsFilter("country,province,-,isp")
	mapper := &MockMapper{
		field:   model.Country,
		replace: "country2",
		matched: true,
	}

	reader := NewFilterReader(database, filter, mapper)
	ast.NotNil(reader)

	ipNet, data, err := reader.LookupNetwork(net.ParseIP("0.0.0.0"))
	ast.Nil(err)
	ast.Equal("0.0.0.0", ipNet.IP.String())
	ast.Equal("country2:province1::isp1", data)

	// calibrate
	calibrate := func(ip net.IP, data map[string]string) (ipNet *net.IPNet, ret map[string]string, updated bool) {
		_, ipNet, err = net.ParseCIDR("0.0.0.1/32")
		ast.Nil(err)
		return ipNet, map[string]string{model.Country: "country1",
			model.Province: "province2",
			model.City:     "city2",
			model.ISP:      "isp2",
		}, true
	}
	reader.Calibrate = calibrate
	ipNet, data, err = reader.LookupNetwork(net.ParseIP("0.0.0.0"))
	ast.Nil(err)
	ast.Equal("0.0.0.1", ipNet.IP.String())
	ast.Equal("country1:province2::isp2", data)

	ast.Equal("country,province,,isp", reader.FieldsCollection())

}

type MockDB struct {
	ipNet *net.IPNet
	data  map[string]string
	err   error
}

func (d *MockDB) LookupNetwork(ip net.IP) (*net.IPNet, map[string]string, error) {
	return d.ipNet, d.data, d.err
}

type MockMapper struct {
	field   string
	replace string
	matched bool
}

func (m *MockMapper) Mapping(field, match string) (string, bool) {
	if field == m.field {
		return m.replace, m.matched
	}
	return "", false
}
