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
	"bytes"
	"encoding/binary"
	"net"
)

// MaxIPv4Uint32 largest possible uint32 value for an IPv4 address
const MaxIPv4Uint32 = 4294967295

var (
	FirstIPv4 = net.IPv4(0, 0, 0, 0)
	LastIPv4  = net.IPv4(255, 255, 255, 255)
	FirstIPv6 = make(net.IP, net.IPv6len)
	LastIPv6  = net.IP{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

// Uint32ToIPv4 uint32
func Uint32ToIPv4(uip uint32) net.IP {
	return net.IPv4(byte(uip>>24&0xFF), byte(uip>>16&0xFF), byte(uip>>8&0xFF), byte(uip&0xFF))
}

// IPv4ToUint32 IPv4 转换为 uint32
func IPv4ToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

// IPv4StrToUint32 IPv4字符串 转换为 uint32 (syntactic sugar)
func IPv4StrToUint32(ipStr string) uint32 {
	if ip := net.ParseIP(ipStr).To4(); ip != nil {
		return IPv4ToUint32(ip)
	}
	return 0
}

// IPToUint32 IP 转换为 uint32, IPv6 返回前 32 位
func IPToUint32(ip net.IP) uint32 {
	if ip4 := ip.To4(); ip4 != nil {
		return IPv4ToUint32(ip4)
	}
	return binary.BigEndian.Uint32(ip.To16()[0:4])
}

// Uint64ToIP uint64 转换为 IP
func Uint64ToIP(uip uint64) net.IP {
	return net.IP{
		byte(uip >> 56 & 0xFF), byte(uip >> 48 & 0xFF), byte(uip >> 40 & 0xFF), byte(uip >> 32 & 0xFF),
		byte(uip >> 24 & 0xFF), byte(uip >> 16 & 0xFF), byte(uip >> 8 & 0xFF), byte(uip & 0xFF),
		0, 0, 0, 0, 0, 0, 0, 0,
	}
}

// Uint64ToIP2 uint64 转换为 IP
func Uint64ToIP2(high, low uint64) net.IP {
	return net.IP{
		byte(high >> 56 & 0xFF), byte(high >> 48 & 0xFF), byte(high >> 40 & 0xFF), byte(high >> 32 & 0xFF),
		byte(high >> 24 & 0xFF), byte(high >> 16 & 0xFF), byte(high >> 8 & 0xFF), byte(high & 0xFF),
		byte(low >> 56 & 0xFF), byte(low >> 48 & 0xFF), byte(low >> 40 & 0xFF), byte(low >> 32 & 0xFF),
		byte(low >> 24 & 0xFF), byte(low >> 16 & 0xFF), byte(low >> 8 & 0xFF), byte(low & 0xFF),
	}
}

// PrevIP returns the IP immediately before the given IP.
func PrevIP(ip net.IP) net.IP {
	res := make(net.IP, len(ip))
	copy(res, ip)
	for i := len(ip) - 1; i >= 0; i-- {
		if res[i] > 0 {
			res[i]--
			break
		}
		res[i] = 0xff
	}
	return res
}

// NextIP returns the IP immediately after the given IP.
func NextIP(ip net.IP) net.IP {
	res := make(net.IP, len(ip))
	copy(res, ip)
	for i := len(ip) - 1; i >= 0; i-- {
		if res[i] < 0xff {
			res[i]++
			break
		}
		res[i] = 0
	}
	return res
}

// IPLess compares two IPs and returns true if the first IP is less than the second.
// IPv4 < IPv6
func IPLess(a, b net.IP) bool {
	if lenA, lenB := len(a), len(b); lenA != lenB {
		return lenA < lenB
	}
	return bytes.Compare(a, b) < 0
}

// IsFirstIP checks if the given IP is the first IP of its kind (IPv4 or IPv6).
func IsFirstIP(ip net.IP, ipv6 bool) bool {
	if ipv6 {
		return ip.Equal(FirstIPv6)
	}
	if len(ip) == net.IPv6len {
		return ip[12] == 0 && ip[13] == 0 && ip[14] == 0 && ip[15] == 0
	}
	return ip.Equal(FirstIPv4)
}

// IsLastIP checks if the given IP is the last IP of its kind (IPv4 or IPv6).
func IsLastIP(ip net.IP, ipv6 bool) bool {
	if ipv6 {
		return ip.Equal(LastIPv6)
	}
	return ip.Equal(LastIPv4)
}
