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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextParser(t *testing.T) {
	ast := assert.New(t)

	ipv4FillResult := func(ip string) []string { return []string{ip} }
	ipv6FillResult := func(ip string) []string { return []string{ip} }

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
