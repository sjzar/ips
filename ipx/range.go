package ipx

import (
	"math/bits"
	"net"
)

// Range IP区间
type Range struct {
	Start net.IP
	End   net.IP
}

// NewRange 初始化IP区间
func NewRange(ipNet *net.IPNet) *Range {
	return &Range{
		Start: ipNet.IP,
		End:   LastIP(ipNet),
	}
}

// Join 合并网段
func (r *Range) Join(r2 *Range) bool {
	// 禁止 起始IP小于原IP区间 或 IP区间不相邻 的IP区间合入
	if IPLess(r2.Start, r.Start) || IPLess(NextIP(r.End), r2.Start) {
		return false
	}

	if !IPLess(r2.End, r.End) {
		r.End = r2.End
	}
	return true
}

// CommonRange 求共同网段
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

// Contains 检查IP区间内是否包含IP
func (r *Range) Contains(ip net.IP) bool {
	return Contains(r.Start, r.End, ip)
}

// JoinIPNet 合并网段
func (r *Range) JoinIPNet(ipNet *net.IPNet) bool {
	// 禁止 起始IP小于原IP区间 或 IP区间不相邻 的IP区间合入
	if IPLess(ipNet.IP, r.Start) || IPLess(NextIP(r.End), ipNet.IP) {
		return false
	}

	if !IPLess(LastIP(ipNet), r.End) {
		r.End = LastIP(ipNet)
	}
	return true
}

// IPNets 输出IP区间对应的CIDR分组
func (r *Range) IPNets() []*net.IPNet {
	start, end := r.Start, r.End
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

// IsEnd IP区间的 End 是否是最后一个IP
func (r *Range) IsEnd() bool {
	return IsZeroIP(NextIP(r.End))
}

// PrefixSameLength 前缀相同长度
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

// SuffixZeroLength 后缀空值长度
func SuffixZeroLength(ip net.IP) int {
	ipLen := len(ip)
	for i := ipLen - 1; i >= 0; i-- {
		if ip[i] != 0 {
			return (ipLen-1-i)*8 + bits.TrailingZeros(uint(ip[i]))
		}
	}
	return ipLen * 8
}

type Ranges []Range

func (r Ranges) Len() int { return len(r) }
func (r Ranges) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r Ranges) Less(i, j int) bool {
	return IPLess(r[i].Start, r[j].Start)
}
