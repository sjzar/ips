package ipx

import (
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestIPv4StrToUint32(t *testing.T) {
	ast := assert.New(t)

	i := IPv4StrToUint32("1.1.1.1")
	ast.Equal("1.1.1.1", Uint32ToIPv4(i).String())
	i = IPv4StrToUint32("fake")
	ast.Equal(uint32(0), i)
}

func TestContains(t *testing.T) {
	ast := assert.New(t)

	start := net.ParseIP("1.1.0.0")
	end := net.ParseIP("1.1.0.255")
	ast.True(Contains(start, end, net.ParseIP("1.1.0.1")))
	ast.True(Contains(start, end, net.ParseIP("1.1.0.0")))
	ast.True(Contains(start, end, net.ParseIP("1.1.0.255")))
	ast.False(Contains(start, end, net.ParseIP("1.1.1.0")))
	ast.True(Contains(net.ParseIP("1.1.0.255"), net.ParseIP("1.1.0.255"), net.ParseIP("1.1.0.255")))

}
