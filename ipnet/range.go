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
	"math/bits"
	"net"
)

// Range represents an IP range with a start and end IP.
type Range struct {
	Start net.IP
	End   net.IP
}

// NewRange initializes an IP range based on the provided IPNet.
// It normalizes the IP addresses to IPv6 length for consistency.
func NewRange(ipNet *net.IPNet) *Range {
	return &Range{
		Start: ipNet.IP.To16(),
		End:   LastIP(ipNet).To16(),
	}
}

// Join attempts to merge the current IP range with another range (r2).
// It forbids merging if the start IP of r2 is before the current range or if the ranges aren't adjacent.
func (r *Range) Join(r2 *Range) bool {

	if IPLess(r2.Start, r.Start) || IPLess(NextIP(r.End), r2.Start) {
		return false
	}

	if !IPLess(r2.End, r.End) {
		r.End = r2.End
	}
	return true
}

// CommonRange finds the common IP range between the current range and another range (r2),
// ensuring that the resulting range contains the specified IP.
func (r *Range) CommonRange(ip net.IP, r2 *Range) bool {
	if IPLess(r2.End, r.Start) || IPLess(NextIP(r.End), r2.Start) {
		return false
	}
	start := r.Start
	if IPLess(start, r2.Start) {
		start = r2.Start
	}
	end := r.End
	if IPLess(r2.End, end) {
		end = r2.End
	}
	if !Contains(start, end, ip) {
		return false
	}
	r.Start, r.End = start, end
	return true
}

// Contains checks if the IP range contains a specific IP.
func (r *Range) Contains(ip net.IP) bool {
	return Contains(r.Start, r.End, ip)
}

// JoinIPNet tries to merge the current IP range with a provided IPNet.
// It forbids merging if the start IP of IPNet is before the current range or if they aren't adjacent.
func (r *Range) JoinIPNet(ipNet *net.IPNet) bool {

	if IPLess(ipNet.IP.To16(), r.Start) || IPLess(NextIP(r.End), ipNet.IP.To16()) {
		return false
	}

	if !IPLess(LastIP(ipNet).To16(), r.End) {
		r.End = LastIP(ipNet).To16()
	}
	return true
}

// IPNets returns the CIDR groups corresponding to the IP range.
func (r *Range) IPNets() []*net.IPNet {
	start, end := r.Start.To16(), r.End.To16()
	bitLength := len(start) * 8
	var result []*net.IPNet
	for {
		cidr := bitLength - SuffixZeroLength(start)
		if cidr < PrefixSameLength(start, end) {
			cidr = PrefixSameLength(start, end)
		}

		ipNet := &net.IPNet{IP: start, Mask: net.CIDRMask(cidr, bitLength)}
		result = append(result, ipNet)
		last := LastIP(ipNet)
		if !IPLess(last, end) {
			return result
		}
		start = NextIP(last)
	}
}

// PrefixSameLength determines the length of the common prefix between two IPs.
func PrefixSameLength(start, end net.IP) int {
	if len(start) != len(end) {
		return 0
	}
	// FIXME 修复越界问题
	endNext := NextIP(end)
	if !start.Equal(end) && start.Equal(endNext) {
		return 0
	}
	end = endNext

	index := 0
	for i := 0; i < len(end); i++ {
		xor := end[i] ^ start[i]
		if xor == 0 {
			index += 8
			continue
		}
		index += bits.LeadingZeros8(xor) + 1
		break
	}
	return index
}

// SuffixZeroLength calculates the length of trailing zeros in the IP.
func SuffixZeroLength(ip net.IP) int {
	ipLen := len(ip)
	for i := ipLen - 1; i >= 0; i-- {
		if ip[i] != 0 {
			return (ipLen-1-i)*8 + bits.TrailingZeros(uint(ip[i]))
		}
	}
	return ipLen * 8
}

// Ranges represents a slice of IP ranges.
type Ranges []Range

// Len returns the number of IP ranges.
func (r Ranges) Len() int { return len(r) }

// Swap swaps the positions of two IP ranges in the slice.
func (r Ranges) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// Less checks if the start IP of the range at index i is less than that at index j.
func (r Ranges) Less(i, j int) bool { return IPLess(r[i].Start, r[j].Start) }
