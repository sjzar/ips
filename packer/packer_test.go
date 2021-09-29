package packer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacker_GetNode(t *testing.T) {
	ast := assert.New(t)

	packer := NewPacker(IPv4, []string{"field1"})
	ast.NotNil(packer)
	ast.Equal(1, len(packer.node))

	// invalid ip
	packer.Load("fake ip", []string{})
	ast.Equal(1, len(packer.node))

	// invalid fields
	packer.Load("0.0.0.0/32", []string{})
	ast.Equal(1, len(packer.node))

	packer.Load("0.0.0.0/32", []string{"field1"})
	ast.Equal(128, len(packer.node)) // 96 + 32 = 128

	// IPv4 to IPv6
	// 0000:0000:0000:0000:0000:FFFF:<ipv4>
	ast.Equal([2]int{80, 0}, packer.node[79])
	ast.Equal([2]int{0, 81}, packer.node[80])

	// cidr conflict
	packer.Load("0.0.0.0/32", []string{"field1"})
	ast.Equal(128, len(packer.node))

	packer.Load("0.0.0.2/32", []string{"field1"})
	ast.Equal(129, len(packer.node))

	packer.Load("255.255.255.255/32", []string{"field1"})
	ast.Equal(160, len(packer.node))            // 128 + 32 = 160
	ast.Equal(8, len(packer.dataChunk.Bytes())) // len(2byte)+field1(2byte)

	_ = packer.Export()
	ast.Equal(160, packer.Meta.NodeCount)
	ast.Equal([]string{"field1"}, packer.Meta.Fields)
	ast.Equal(IPv4, packer.Meta.IPVersion)
	ast.Equal(1296, packer.Meta.TotalSize) // nodeChunk(160 * 8 + loopNode(8byte)) + dataChunk(8byte)
}
