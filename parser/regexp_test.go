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

func TestIPv4Regexp(t *testing.T) {
	ast := assert.New(t)

	type instance struct {
		str  string
		find int
	}

	instances := []instance{
		{str: "", find: 0},
		{str: "1.1.1.1", find: 1},
		{str: "0.1.1.1", find: 1},
		{str: "1.1.1.01", find: 1}, // 1.1.1.0
		{str: "hello 1.1.1.1 ips 6.2.1.255", find: 2},
		{str: "1.2.3.4 266.1.1.2", find: 2}, // 1.2.3.4 66.1.1.2
		{str: "....2.3.3.4", find: 1},
		{str: "666123.123.1.25", find: 1},
		{str: "fe80::1 6.7.8.9", find: 1},
	}

	for index, inst := range instances {
		ast.Equal(inst.find, len(IPv4Regexp.FindAllStringIndex(inst.str, -1)), "index: %d str: %s", index, inst.str)
	}
}

func TestIPv6Regexp(t *testing.T) {
	ast := assert.New(t)

	type instance struct {
		str  string
		find int
	}

	instances := []instance{
		{str: "", find: 0},
		{str: "fe80::1", find: 1},
		{str: "fe80::1 fe80::2", find: 2},
		{str: "1:2:3:4:5:6:7:8", find: 1},
		{str: "1:2:3:4:5:6::8 1:2:3:4:5::7:8 1:2:3:4::6:7:8 1:2:3::5:6:7:8 1:2::4:5:6:7:8 1::3:4:5:6:7:8", find: 6},
		{str: "abc::111", find: 1}, // bc::111
		{str: "abc:111", find: 0},
		{str: "1:::::::::1", find: 4}, // 1:: / :: / :: / :(error) / ::1
		{str: "::FFFF:1.1.1.1", find: 1},
		{str: "::1.1.1.1", find: 1},
		{str: "fe80::1%eth0", find: 1},
	}

	for index, inst := range instances {
		ast.Equal(inst.find, len(IPv6Regexp.FindAllStringIndex(inst.str, -1)), "index: %d str: %s", index, inst.str)
	}
}
