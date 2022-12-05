package ipx

import (
	"bytes"
	"net"
)

const MaxIPv4Uint32 = 4294967295

var (
	ZeroIPv4 = make(net.IP, net.IPv4len)
	ZeroIPv6 = make(net.IP, net.IPv6len)
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

// IsZeroIP 是否是起始IP
func IsZeroIP(ip net.IP) bool {
	return ip.Equal(ZeroIPv4) || ip.Equal(ZeroIPv6)
}
