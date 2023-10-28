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
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint32ToIPv4(t *testing.T) {
	ast := assert.New(t)

	ip := Uint32ToIPv4(128)
	ast.Equal("0.0.0.128", ip.String())

	for times := 0; times < 10; times++ {
		i := rand.Uint32()
		ast.Equal(i, IPv4ToUint32(Uint32ToIPv4(i)))
	}
}

func TestIPv4StrToUint32(t *testing.T) {
	ast := assert.New(t)

	i := IPv4StrToUint32("1.1.1.1")
	ast.Equal("1.1.1.1", Uint32ToIPv4(i).String())
	i = IPv4StrToUint32("fake")
	ast.Equal(uint32(0), i)
}

func TestPrevIP(t *testing.T) {
	tests := []struct {
		input    net.IP
		expected net.IP
	}{
		{net.ParseIP("192.168.1.1"), net.ParseIP("192.168.1.0")},
		{net.ParseIP("192.168.1.0"), net.ParseIP("192.168.0.255")},
		{net.ParseIP("0.0.0.0"), net.ParseIP("::fffe:ffff:ffff")},
		{net.ParseIP("::1"), net.ParseIP("::")},
		{net.ParseIP("::"), net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")},
	}

	for _, test := range tests {
		output := PrevIP(test.input)
		if !output.Equal(test.expected) {
			t.Errorf("For input %v, expected %v, but got %v", test.input, test.expected, output)
		}
	}
}

func TestNextIP(t *testing.T) {
	tests := []struct {
		input    net.IP
		expected net.IP
	}{
		{net.ParseIP("192.168.1.0"), net.ParseIP("192.168.1.1")},
		{net.ParseIP("192.168.0.255"), net.ParseIP("192.168.1.0")},
		{net.ParseIP("255.255.255.255"), net.ParseIP("::1:0:0:0")},
		{net.ParseIP("::"), net.ParseIP("::1")},
		{net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"), net.ParseIP("::")},
	}

	for _, test := range tests {
		output := NextIP(test.input)
		if !output.Equal(test.expected) {
			t.Errorf("For input %v, expected %v, but got %v", test.input, test.expected, output)
		}
	}
}
