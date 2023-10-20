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

package ipnet

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	ast := assert.New(t)

	_, ipNet, err := net.ParseCIDR("0.0.0.0/32")
	ast.Nil(err)

	// new Range
	ipr := NewRange(ipNet)
	ast.Equal("0.0.0.0", ipr.Start.String())
	ast.Equal("0.0.0.0", ipr.End.String())

	// join network
	_, ipNet, err = net.ParseCIDR("0.0.0.1/32")
	ast.Nil(err)
	_ipr := NewRange(ipNet)
	ok := ipr.Join(_ipr)
	ast.True(ok)
	ast.Equal("0.0.0.1", ipr.End.String())

	// join wrong network
	_, ipNet, err = net.ParseCIDR("1.0.0.0/32")
	ast.Nil(err)
	_ipr = NewRange(ipNet)
	ok = ipr.Join(_ipr)
	ast.False(ok)
	ast.Equal("0.0.0.1", ipr.End.String())

	// output
	ipNets := ipr.IPNets()
	ast.Equal(1, len(ipNets))
	ast.Equal("0.0.0.0/31", ipNets[0].String())

	ipr.End = net.ParseIP("0.0.0.255").To4()
	ipNets = ipr.IPNets()
	ast.Equal(1, len(ipNets))
	ast.Equal("0.0.0.0/24", ipNets[0].String())

	ipr.End = net.ParseIP("0.0.4.255").To4()
	ipNets = ipr.IPNets()
	ast.Equal(2, len(ipNets))

	// join sub network
	_, ipNet, err = net.ParseCIDR("58.82.200.0/23")
	ast.Nil(err)
	ipr = NewRange(ipNet)
	ast.Equal("58.82.200.0", ipr.Start.String())
	ast.Equal("58.82.201.255", ipr.End.String())
	_, ipNet, err = net.ParseCIDR("58.82.203.0/21")
	ast.Nil(err)
	ok = ipr.JoinIPNet(ipNet)
	ast.True(ok)
	ast.Equal("58.82.207.255", ipr.End.String())

	_, ipNet, err = net.ParseCIDR("58.82.203.0/20")
	ast.Nil(err)
	ok = ipr.JoinIPNet(ipNet)
	ast.False(ok)
}

func TestRange_CommonRange(t *testing.T) {
	ast := assert.New(t)

	_, ipNet1, err := net.ParseCIDR("58.82.203.0/21")
	ast.Nil(err)
	ipr1 := NewRange(ipNet1)
	ast.Equal("58.82.200.0", ipr1.Start.String())
	ast.Equal("58.82.207.255", ipr1.End.String())

	_, ipNet2, err := net.ParseCIDR("58.82.200.0/23")
	ast.Nil(err)
	ipr2 := NewRange(ipNet2)
	ast.Equal("58.82.200.0", ipr2.Start.String())
	ast.Equal("58.82.201.255", ipr2.End.String())

	// Example1
	// IP:        *
	// ipr1: [               ]
	// ipr2: [      ]
	// Result:   [      ]
	ast.True(ipr1.CommonRange(ipr1.Start, ipr2))
	ast.Equal("58.82.200.0", ipr1.Start.String())
	ast.Equal("58.82.201.255", ipr1.End.String())

	_, ipNet3, err := net.ParseCIDR("58.82.201.0/24")
	ast.Nil(err)
	ipr3 := NewRange(ipNet3)
	ast.Equal("58.82.201.0", ipr3.Start.String())
	ast.Equal("58.82.201.255", ipr3.End.String())

	// Example2
	// IP:        *
	// ipr1: [               ]
	// ipr3:          [      ]
	// Result: Error ipr2 is not include IP
	ast.False(ipr1.CommonRange(ipr1.Start, ipr3))

	// Example3
	// IP:                 *
	// ipr1: [               ]
	// ipr2:          [      ]
	// Result:            [      ]
	ast.True(ipr1.CommonRange(ipr3.Start, ipr3))
	ast.Equal("58.82.201.0", ipr1.Start.String())
	ast.Equal("58.82.201.255", ipr1.End.String())

	_, ipNet4, err := net.ParseCIDR("58.82.208.0/24")
	ast.Nil(err)
	ipr4 := NewRange(ipNet4)
	ast.Equal("58.82.208.0", ipr4.Start.String())
	ast.Equal("58.82.208.255", ipr4.End.String())

	// Example4
	// IP:        *
	// ipr1: [      ]
	// ipr4:          [      ]
	// Result: Error ipr4 is not adjacent
	ast.False(ipr1.CommonRange(ipr1.Start, ipr4))
}
