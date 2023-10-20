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
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
