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
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskLess(t *testing.T) {
	ast := assert.New(t)

	_, ipNet1, err := net.ParseCIDR("0.0.0.0/16")
	ast.Nil(err)

	_, ipNet2, err := net.ParseCIDR("0.0.0.0/24")
	ast.Nil(err)

	ast.False(MaskLess(ipNet1, ipNet2))
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
