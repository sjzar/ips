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
)

// MaskLess compares two CIDRs.
// Returns true if the first CIDR has a larger prefix length (smaller network) than the second.
// IPv4 networks are considered smaller than IPv6.
func MaskLess(a, b *net.IPNet) bool {
	if lenA, lenB := len(a.IP), len(b.IP); lenA != lenB {
		return lenA < lenB
	}
	aOnes, _ := a.Mask.Size()
	bOnes, _ := b.Mask.Size()
	return aOnes > bOnes
}

// LastIP returns the last IP address within the provided IPNet.
func LastIP(ipNet *net.IPNet) net.IP {
	ip, mask := ipNet.IP, ipNet.Mask
	ipLen := len(ip)
	res := make(net.IP, ipLen)
	for i := 0; i < ipLen; i++ {
		res[i] = ip[i] | ^mask[i]
	}
	return res
}

// Contains checks if the given IP is within the range of start to end (inclusive of start and exclusive of end).
func Contains(start, end, ip net.IP) bool {
	return !IPLess(ip, start) && IPLess(ip, NextIP(end))
}
