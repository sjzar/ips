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

package ipx

import (
	"bytes"
	"net"
)

const MaxIPv4Uint32 = 4294967295

var (
	FirstIPv4 = net.IPv4(0, 0, 0, 0)
	LastIPv4  = net.IPv4(255, 255, 255, 255)
	FirstIPv6 = make(net.IP, net.IPv6len)
	LastIPv6  = net.IP{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

// Uint32ToIPv4 uint32 转换为 IP
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

// IPNetMaskLess CIDR大小比较
// IPv4 < IPv6
// ones越大，子网范围越小
func IPNetMaskLess(a, b *net.IPNet) bool {
	if lenA, lenB := len(a.IP), len(b.IP); lenA != lenB {
		return lenA < lenB
	}
	aOnes, _ := a.Mask.Size()
	bOnes, _ := b.Mask.Size()
	return aOnes > bOnes
}

// PrevIP 上一个IP
func PrevIP(ip net.IP) net.IP {
	res := make(net.IP, len(ip))
	for i := len(ip) - 1; i >= 0; i-- {
		res[i] = ip[i] - 1
		if res[i] != 0xff {
			copy(res, ip[0:i])
			break
		}
	}
	return res
}

// NextIP 下一个IP
func NextIP(ip net.IP) net.IP {
	res := make(net.IP, len(ip))
	for i := len(ip) - 1; i >= 0; i-- {
		res[i] = ip[i] + 1
		if res[i] != 0 {
			copy(res, ip[0:i])
			break
		}
	}
	return res
}

// LastIP IPNet最后一个IP
func LastIP(ipNet *net.IPNet) net.IP {
	ip, mask := ipNet.IP, ipNet.Mask
	ipLen := len(ip)
	res := make(net.IP, ipLen)
	for i := 0; i < ipLen; i++ {
		res[i] = ip[i] | ^mask[i]
	}
	return res
}

// IPLess IP大小比较
// IPv4 < IPv6
func IPLess(a, b net.IP) bool {
	if lenA, lenB := len(a), len(b); lenA != lenB {
		return lenA < lenB
	}
	return bytes.Compare(a, b) < 0
}

// Contains 检查IP区间内是否包含IP
func Contains(start, end, ip net.IP) bool {
	return !IPLess(ip, start) && IPLess(ip, NextIP(end))
}

// IsFirstIP 是否是起始IP
func IsFirstIP(ip net.IP, ipv6 bool) bool {
	if ipv6 {
		return ip.Equal(FirstIPv6)
	}
	if len(ip) == net.IPv6len {
		return ip[12] == 0 && ip[13] == 0 && ip[14] == 0 && ip[15] == 0
	}
	return ip.Equal(FirstIPv4)
}

// IsLastIP 是否是最后一个IP
func IsLastIP(ip net.IP, ipv6 bool) bool {
	if ipv6 {
		return ip.Equal(LastIPv6)
	}
	return ip.Equal(LastIPv4)
}
