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

package ipdb

import (
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

func TestWriter(t *testing.T) {
	ast := assert.New(t)

	meta := model.Meta{
		IPVersion: model.IPv4,
		Fields:    []string{"field1"},
	}

	writer := NewWriter(meta, nil)
	ast.NotNil(writer)
	ast.Equal(1, len(writer.node))

	// invalid fields
	_, ipNet, err := net.ParseCIDR("0.0.0.0/32")
	ast.Nil(err)
	err = writer.Insert(ipx.NewRange(ipNet), []string{})
	ast.Equal(errors.ErrInvalidFieldsLength, err)

	err = writer.Insert(ipx.NewRange(ipNet), []string{"field1"})
	ast.Nil(err)
	ast.Equal(128, len(writer.node)) // 96 + 32 = 128

	// IPv4 to IPv6
	// 0000:0000:0000:0000:0000:FFFF:<ipv4>
	ast.Equal([2]int{80, 0}, writer.node[79])
	ast.Equal([2]int{0, 81}, writer.node[80])

	// cidr overwrite
	err = writer.Insert(ipx.NewRange(ipNet), []string{"field1"})
	ast.Nil(err)
	ast.Equal(128, len(writer.node))

	_, ipNet, err = net.ParseCIDR("0.0.0.2/32")
	ast.Nil(err)
	err = writer.Insert(ipx.NewRange(ipNet), []string{"field1"})
	ast.Nil(err)

	_, ipNet, err = net.ParseCIDR("255.255.255.255/32")
	ast.Nil(err)
	err = writer.Insert(ipx.NewRange(ipNet), []string{"field1"})
	ast.Nil(err)
	ast.Equal(160, len(writer.node))            // 128 + 32 = 160
	ast.Equal(8, len(writer.dataChunk.Bytes())) // len(2byte)+field1(2byte)

	_ = writer.Save(io.Discard)
	ast.Equal(160, writer.Meta.NodeCount)
	ast.Equal([]string{"field1"}, writer.Meta.Fields)
	ast.Equal(model.IPv4, writer.Meta.IPVersion)
	ast.Equal(1296, writer.Meta.TotalSize) // nodeChunk(160 * 8 + loopNode(8byte)) + dataChunk(8byte)
}
