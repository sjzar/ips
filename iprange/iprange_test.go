package iprange

import (
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPRange(t *testing.T) {
	ast := assert.New(t)

	_, ipNet, err := net.ParseCIDR("0.0.0.0/32")
	ast.Nil(err)

	// new IPRange
	ipRange := NewIPRange(ipNet)
	ast.Equal("0.0.0.0", ipRange.Start.String())
	ast.Equal("0.0.0.0", ipRange.End.String())

	// join network
	_, ipNet, err = net.ParseCIDR("0.0.0.1/32")
	ast.Nil(err)
	ipRange.Join(ipNet)
	ast.Equal("0.0.0.1", ipRange.End.String())

	// join wrong network
	_, ipNet, err = net.ParseCIDR("1.0.0.0/32")
	ast.Nil(err)
	ipRange.Join(ipNet)
	ast.Equal("0.0.0.1", ipRange.End.String())

	// output
	ipNets := ipRange.IPNets()
	ast.Equal(1, len(ipNets))
	ast.Equal("0.0.0.0/31", ipNets[0].String())

	ipRange.End = net.ParseIP("0.0.0.255").To4()
	ipNets = ipRange.IPNets()
	ast.Equal(1, len(ipNets))
	ast.Equal("0.0.0.0/24", ipNets[0].String())

	ipRange.End = net.ParseIP("0.0.4.255").To4()
	ipNets = ipRange.IPNets()
	ast.Equal(2, len(ipNets))
}

func TestIPNetMaskLess(t *testing.T) {
	ast := assert.New(t)

	_, ipNet1, err := net.ParseCIDR("0.0.0.0/16")
	ast.Nil(err)

	_, ipNet2, err := net.ParseCIDR("0.0.0.0/24")
	ast.Nil(err)

	ast.False(IPNetMaskLess(ipNet1, ipNet2))
}

func TestUint32ToIPv4(t *testing.T) {
	ast := assert.New(t)

	ip := Uint32ToIPv4(128)
	ast.Equal("0.0.0.128", ip.String())

	for times := 0; times < 10; times++ {
		i := rand.Uint32()
		ast.Equal(i, IPv4ToUint32(Uint32ToIPv4(i)))
	}
}
