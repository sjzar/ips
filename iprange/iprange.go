package iprange

import (
	"bytes"
	"math/bits"
	"net"
)

const MaxIPv4Uint32 = 4294967295

// IPRange IP区间
type IPRange struct {
	Start net.IP
	End   net.IP
}

// NewIPRange 初始化IP区间
func NewIPRange(ipNet *net.IPNet) *IPRange {
	return &IPRange{
		Start: ipNet.IP,
		End:   LastIP(ipNet),
	}
}

// Join 合并网段
func (r *IPRange) Join(ipNet *net.IPNet) {
	if NextIP(r.End).Equal(ipNet.IP) {
		r.End = LastIP(ipNet)
	}
}

// IPNets 输出IP区间对应的CIDR分组
func (r *IPRange) IPNets() []*net.IPNet {
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

// IPLess IP大小比较
// IPv4 < IPv6
func IPLess(a, b net.IP) bool {
	if lenA, lenB := len(a), len(b); lenA != lenB {
		return lenA < lenB
	}
	return bytes.Compare(a, b) < 0
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

type IPRanges []IPRange

func (r IPRanges) Len() int { return len(r) }
func (r IPRanges) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r IPRanges) Less(i, j int) bool {
	return IPLess(r[i].Start, r[j].Start)
}

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
