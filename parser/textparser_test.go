package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextParser(t *testing.T) {
	ast := assert.New(t)

	ipv4FillResult := func(ip string) string { return ip }
	ipv6FillResult := func(ip string) string { return ip }

	type instance struct {
		str    string
		result string
	}

	instances := []instance{
		{str: "", result: ""},
		{str: "123", result: "123"},
		{str: "1.1.1.1", result: "1.1.1.1 [1.1.1.1] "},
		{str: "abc1.1.1.1def", result: "abc1.1.1.1 [1.1.1.1] def"},
		{str: "666.123.231.321g99::80z", result: "666.123.231.32 [66.123.231.32] 1g99::80 [99::80] z"},
		{str: "::ffff:1.2.3.4.5.6.7.8", result: "::ffff:1.2.3.4 [::ffff:1.2.3.4] .5.6.7.8 [5.6.7.8] "},
	}

	for index, inst := range instances {
		parser := NewTextParser(inst.str)
		parser.IPv4FillResult = ipv4FillResult
		parser.IPv6FillResult = ipv6FillResult
		ast.Equal(inst.result, parser.Parse().String(), "index: %d str: %s", index, inst.str)
	}
}
