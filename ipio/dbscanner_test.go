/*
 * Copyright (c) 2022 shenjunzheng@gmail.com
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

package ipio

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
	"github.com/sjzar/ips/rewriter"
)

func TestReader(t *testing.T) {
	ast := assert.New(t)

	_, ipNet, err := net.ParseCIDR("0.0.0.0/32")
	ast.Nil(err)
	ipr := ipx.NewRange(ipNet)

	database := &MockDB{
		ipRange: ipr,
		data: map[string]string{
			model.Country:  "country1",
			model.Province: "province1",
			model.City:     "city1",
			model.ISP:      "isp1",
		},
	}
	selector := NewFieldSelector("country,province,,isp")
	dl := rewriter.NewDataLoader()
	dl.LoadString("isp\tisp1\tisp2")
	rw := rewriter.NewDataRewriter(dl, nil)

	reader := NewDBScanner(database, selector, rw)
	ast.NotNil(reader)

	ipr, values, err := reader.Find(net.ParseIP("0.0.0.0"))
	ast.Nil(err)
	ast.Equal("0.0.0.0", ipNet.IP.String())
	ast.Equal([]string{"country1", "province1", "", "isp2"}, values)
	ast.Equal([]string{"country", "province", "", "isp"}, reader.Meta().Fields)
}

func TestReader_Scan(t *testing.T) {
	ast := assert.New(t)

	database := &MockDB2{
		mask:   2,
		fields: []string{"country", "isp"},
		data:   map[string]string{"country": "中国", "isp": "电信"},
	}
	selector := NewFieldSelector("country,isp")
	reader := NewDBScanner(database, selector, nil)
	ast.NotNil(reader)

	err := reader.Init(model.Meta{})
	ast.Nil(err)
	ok := reader.Scan()
	ast.True(ok)
	ipr, values := reader.Result()
	ast.Equal("0.0.0.0/0", ipr.IPNets()[0].String())
	ast.Equal([]string{"中国", "电信"}, values)

	database.diffData = true
	err = reader.Init(model.Meta{IPVersion: model.IPv4})
	ast.Nil(err)
	ok = reader.Scan()
	ast.True(ok)
	ipr, values = reader.Result()
	ast.NotNil(values)
	ast.Equal("0.0.0.0/2", ipr.IPNets()[0].String())
	ok = reader.Scan()
	ast.True(ok)
	ipr, values = reader.Result()
	ast.Equal("64.0.0.0/2", ipr.IPNets()[0].String())
	ok = reader.Scan()
	ast.True(ok)
	ipr, values = reader.Result()
	ast.Equal("128.0.0.0/2", ipr.IPNets()[0].String())
	ok = reader.Scan()
	ast.True(ok)
	ipr, values = reader.Result()
	ast.Equal("192.0.0.0/2", ipr.IPNets()[0].String())
	ok = reader.Scan()
	ast.False(ok)
	ast.Nil(reader.Err())

	err = reader.Init(model.Meta{IPVersion: model.IPv4 | model.IPv6})
	ast.Nil(err)
	for i := 0; i < 5; i++ {
		ok = reader.Scan()
		ast.True(ok)
	}
	ipr, values = reader.Result()
	ast.Equal("::/2", ipr.IPNets()[0].String())

	err = reader.Init(model.Meta{IPVersion: model.IPv6})
	ast.Nil(err)
	ok = reader.Scan()
	ast.True(ok)
	ipr, values = reader.Result()
	ast.Equal("::/2", ipr.IPNets()[0].String())
	ok = reader.Scan()
	ast.True(ok)
	ipr, values = reader.Result()
	ast.Equal("4000::/2", ipr.IPNets()[0].String())
}

type MockDB struct {
	ipRange *ipx.Range
	data    map[string]string
	err     error
}

func (d *MockDB) Find(ip net.IP) (*ipx.Range, map[string]string, error) {
	return d.ipRange, d.data, d.err
}

func (d *MockDB) Meta() model.Meta {
	return model.Meta{
		IPVersion: model.IPv4,
		Fields:    []string{"country", "province", "city", "isp"},
	}
}

func (d *MockDB) Close() error {
	return nil
}

type MockDB2 struct {
	mask     int
	fields   []string
	data     map[string]string
	diffData bool
}

func (d *MockDB2) Find(ip net.IP) (*ipx.Range, map[string]string, error) {
	_, ipNet, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip.String(), d.mask))
	ipr := ipx.NewRange(ipNet)
	if d.diffData {
		return ipr, map[string]string{"country": ip.String()}, nil
	}
	return ipr, d.data, nil
}

func (d *MockDB2) Meta() model.Meta {
	return model.Meta{
		IPVersion: model.IPv4 | model.IPv6,
		Fields:    d.fields,
	}
}

func (d *MockDB2) Close() error {
	return nil
}
