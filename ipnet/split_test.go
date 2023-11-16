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
	"fmt"
	"math/big"
	"net"
	"testing"
)

// TestIPToBigInt tests the conversion of an IP address to a big.Int.
func TestIPToBigInt(t *testing.T) {
	testCases := []struct {
		name     string
		ip       net.IP
		expected *big.Int
	}{
		{
			name:     "IPv4 address",
			ip:       net.ParseIP("192.0.2.1"),
			expected: big.NewInt(0).SetBytes([]byte{0xff, 0xff, 192, 0, 2, 1}),
		},
		{
			name:     "IPv6 address",
			ip:       net.ParseIP("2001:db8::1"),
			expected: big.NewInt(0).SetBytes(net.ParseIP("2001:db8::1")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IPToBigInt(tc.ip)
			if result.Cmp(tc.expected) != 0 {
				t.Errorf("IPToBigInt(%v) = %v, want %v", tc.ip, result, tc.expected)
			}
		})
	}
}

// TestBigIntToIP tests the conversion of a big.Int to an IP address.
func TestBigIntToIP(t *testing.T) {
	testCases := []struct {
		name     string
		intValue *big.Int
		expected net.IP
	}{
		{
			name:     "IPv4 address",
			intValue: big.NewInt(0).SetBytes(net.ParseIP("192.0.2.1")),
			expected: net.ParseIP("192.0.2.1"),
		},
		{
			name:     "IPv6 address",
			intValue: big.NewInt(0).SetBytes(net.ParseIP("2001:db8::1")),
			expected: net.ParseIP("2001:db8::1"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := BigIntToIP(tc.intValue)
			if !result.Equal(tc.expected) {
				t.Errorf("BigIntToIP(%v) = %v, want %v", tc.intValue, result, tc.expected)
			}
		})
	}
}

func TestSplitIPNet(t *testing.T) {
	testCases := []struct {
		name   string
		start  net.IP
		end    net.IP
		num    int
		expect []net.IP
	}{
		{
			name:  "IPv4 range step 0",
			start: net.ParseIP("192.0.2.0"),
			end:   net.ParseIP("192.0.2.255"),
			num:   4,
			expect: []net.IP{
				net.ParseIP("192.0.2.0"),
				net.ParseIP("192.0.2.85"),
				net.ParseIP("192.0.2.170"),
				net.ParseIP("192.0.2.255"),
			},
		},
		{
			name:  "IPv4 range step > 1",
			start: net.ParseIP("0.0.0.0"),
			end:   net.ParseIP("255.255.255.255"),
			num:   10,
			expect: []net.IP{
				net.ParseIP("0.0.0.0"),
				net.ParseIP("45.213.238.0"),
				net.ParseIP("65.64.182.200"),
				net.ParseIP("74.85.97.120"),
				net.ParseIP("85.184.194.91"),
				net.ParseIP("99.64.87.96"),
				net.ParseIP("137.148.164.152"),
				net.ParseIP("184.107.61.72"),
				net.ParseIP("206.18.212.16"),
				net.ParseIP("255.255.255.255"),
			},
		},
		{
			name:  "IPv4 range step [0,1]",
			start: net.ParseIP("180.0.0.0"),
			end:   net.ParseIP("185.0.0.0"),
			num:   5,
			expect: []net.IP{
				net.ParseIP("180.0.0.0"),
				net.ParseIP("180.240.102.0"),
				net.ParseIP("181.224.204.0"),
				net.ParseIP("183.38.4.164"),
				net.ParseIP("185.0.0.0"),
			},
		},
		{
			name:  "IPv6 range step 0",
			start: net.ParseIP("2001:db8::"),
			end:   net.ParseIP("2001:db8::ffff"),
			num:   4,
			expect: []net.IP{
				net.ParseIP("2001:db8::"),
				net.ParseIP("2001:db8::5555"),
				net.ParseIP("2001:db8::aaaa"),
				net.ParseIP("2001:db8::ffff"),
			},
		},
		{
			name:  "IPv6 range step > 1",
			start: net.ParseIP("::"),
			end:   net.ParseIP("ffff::"),
			num:   10,
			expect: []net.IP{
				net.ParseIP("::"),
				net.ParseIP("::568b:af60"),
				net.ParseIP("::ca82:be20"),
				net.ParseIP("80.194.168.0"),
				net.ParseIP("193.227.132.0"),
				net.ParseIP("2001:0:4c25:df20::"),
				net.ParseIP("2001:0:bbd0:6200::"),
				net.ParseIP("2002:4721:1400::"),
				net.ParseIP("2002:ae3f:127c::"),
				net.ParseIP("ffff::"),
			},
		},
		{
			name:  "IPv6 range step [0,1]",
			start: net.ParseIP("::"),
			end:   net.ParseIP("::2d01:ffff"),
			num:   4,
			expect: []net.IP{
				net.ParseIP("::"),
				net.ParseIP("::1029:6aaa"),
				net.ParseIP("::187a:b97"),
				net.ParseIP("::2d01:ffff"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SplitIPNet(tc.start, tc.end, tc.num)
			if len(result) != len(tc.expect) {
				t.Fatalf("SplitIPNet(%v, %v, %d) returned %d IPs, want %d", tc.start, tc.end, tc.num, len(result), len(tc.expect))
			}
			for i := range result {
				if !result[i].Equal(tc.expect[i]) {
					t.Errorf("SplitIPNet(%v, %v, %d)[%d] = %v, want %v", tc.start, tc.end, tc.num, i, result[i], tc.expect[i])
				}
			}
		})
	}
}

func TestIPNet(t *testing.T) {

	start := net.ParseIP("::0")
	end := net.ParseIP("::18f2:a2c8")
	fmt.Printf("%#v\n", end)
	// splitIPNetIPv6(start, end, 100)
	ret := splitIPNetIPv6(start, end, 100)

	for _, ip := range ret {
		fmt.Printf("%s\n", ip)
	}

}

func TestGetIndex(t *testing.T) {

	ip := net.ParseIP("255.255.255.255")
	index := GetIndex(BaseIPv4, ip)
	fmt.Println(ip, index, BaseIPv4[index])

}
