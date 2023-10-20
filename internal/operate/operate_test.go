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

package operate

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjzar/ips/ipnet"
	"github.com/sjzar/ips/pkg/model"
)

func TestIPOperateChain(t *testing.T) {
	ast := assert.New(t)

	i := &model.IPInfo{IP: net.ParseIP("0.0.0.0")}

	chain := &IPOperateChain{}
	chain.Use(func(info *model.IPInfo) error {
		info.IP = ipnet.NextIP(info.IP)
		return nil
	})
	chain.Use(func(info *model.IPInfo) error {
		info.IP = ipnet.NextIP(info.IP)
		return nil
	})

	err := chain.Do(i)
	ast.Nil(err)

	ast.Equal(net.ParseIP("0.0.0.2"), i.IP)

	defer func() {
		if r := recover(); r != nil {
			ast.Equal("IPOperateChain is nil", r)
		}
	}()

	var chain2 *IPOperateChain
	chain2.Use(func(info *model.IPInfo) error {
		info.IP = ipnet.NextIP(info.IP)
		return nil
	})
}
