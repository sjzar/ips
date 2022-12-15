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
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

/* IPDB Format
	+--------------------------------+--------------------------------+
	|    MetaData Length (4byte)     |     MetaData (Json Format)     |
	+--------------------------------+--------------------------------+
	|                 Node Chunk (Prefix Tree / Trie)                 |
	+--------------------------------+--------------------------------+
	|                            Data Chunk                           |
	+--------------------------------+--------------------------------+

* 文件被划分为3块，分别是MetaData、NodeChunk、DataChunk
* MetaData是一个json，保存了数据库的构建信息和查询所需的一些元数据
* NodeChunk是前缀树(Trie)，每个Node为8byte，用于保存下一跳offset位置信息，若offset大于Node数量，将跳转到DataChunk，表示查询到了结果
* DataChunk用于保存IP库数据，相同的数据仅保存一份，减少数据冗余

# metadata 示例
{
	"build": 1632971142,    // 构建时间
	"ip_version": 1,        // IP库版本 IPv4:0x1 IPv6:0x2
	"languages": {
		"CN": 0             // 语言 & 字段偏移量
	},
	"node_count": 8705098,  // Node数量
	"total_size": 90028407, // NodeChunk+DataChunk 数据大小
	"fields": [             // DataChunk中每组数据的字段
		"country_name",
		"region_name",
		"city_name",
		"owner_domain",
		"isp_domain",
		"latitude",
		"longitude",
		"timezone",
		"utc_offset",
		"china_admin_code",
		"idd_code",
		"country_code",
		"continent_code"
	]
}

	+--------------------------------+--------------------------------+--------------------------------+
	| Data Length (2byte) | Data Fields(<country>\t<province>\t<city>\t<isp>\t<country>\t<province>)   |
	+--------------------------------+--------------------------------+--------------------------------+
	| Data Length (2byte) | Data Fields(<country>\t<province>\t<city>\t<isp>\t<country>\t<province>)   |
	+--------------------------------+--------------------------------+--------------------------------+
* DataChunk中，数据被分为两个部分，第一部分是长度，第二部分是数据
* 数据部分使用`\t`分隔字段，多语言版本的IP库采用字段offset返回不同语言的数据

* 查询
* CIDR地址是IP+子网掩码的一种网段的描述方式，例如10.0.0.1/8，表示子网掩码为8位（255.0.0.0），CIDR中所有的IP，子网掩码部分完全相同
* 将IP地址看作32bit/128bit的二进制字符串，在NodeChunk中使用前缀树(Trie)从前往后查询，匹配到CIDR后，跳转到DataChunk中返回对应数据
* CIDR分组越多，Node数量越大
* CIDR禁止互相包含，例如将 10.0.0.1/8 和 10.0.0.1/16 设置为不同的结果，在查询过程中，由于先匹配到了10.0.0.1/8，所以10.0.0.1/16 不生效，需要做拆分
* 详细的查询过程，可以参考这篇论文 IPv4 route lookup on Linux: https://vincent.bernat.ch/en/blog/2017-ipv4-route-lookup-linux

* 打包
* 构建前缀树，将不同的数据按照Load顺序放入DataChunk中
* 前缀树按照IPv6的规格进行构建，IPv4需要补全前面96bit的子网数据，80bit的0和16bit的1，对应IPv6的映射地址（::FFFF:<IPv4>）；查询IPv4时，支持快速offset到96bit mask位置起进行查询
* ipdb格式数据库理论上支持IPv4/IPv6数据处于同一个数据库中，但是需要注意::FFFF:路径上禁止存在其他CIDR记录，保证到IPv4查询的路径畅通（修改查询SDK可以更好的支持IPv4/IPv6同时查询）

*/

const (
	FieldsSep = "\t"
)

// Writer IPDB 写入工具
type Writer struct {
	Meta      model.IPDBMeta
	node      [][2]int
	dataHash  map[string]int
	dataChunk *bytes.Buffer
}

// NewWriter 初始化 IPDB 写入实例
func NewWriter(meta model.Meta, languages map[string]int) *Writer {
	if len(languages) == 0 {
		languages = map[string]int{"CN": 0}
	}

	return &Writer{
		Meta: model.IPDBMeta{
			Build:     int(time.Now().Unix()),
			IPVersion: meta.IPVersion,
			Languages: languages,
			Fields:    model.FieldsReplace(CommonFieldsMap, meta.Fields),
		},
		node:      [][2]int{{}},
		dataChunk: &bytes.Buffer{},
		dataHash:  make(map[string]int),
	}
}

// Insert 插入数据
func (p *Writer) Insert(ipr *ipx.Range, values []string) error {
	if len(values) != len(p.Meta.Fields) {
		return errors.ErrInvalidFieldsLength
	}

	for _, ipNet := range ipr.IPNets() {
		if err := p.insert(ipNet, values); err != nil {
			return err
		}
	}

	return nil
}

// Save 保存数据
func (p *Writer) Save(w io.Writer) error {

	// Node Chunk
	nodeChunk := &bytes.Buffer{}
	p.Meta.NodeCount = len(p.node)
	for i := 0; i < p.Meta.NodeCount; i++ {
		for j := 0; j < 2; j++ {
			// 小于0: 数据记录，设置为NodeLength+Data偏移量
			// 等于0: 空值，设置为NodeLength
			// 大于0: Node跳转记录，不做调整
			if p.node[i][j] <= 0 {
				p.node[i][j] = p.Meta.NodeCount - p.node[i][j]
			}
			nodeChunk.Write(IntToBinaryBE(p.node[i][j], 32))
		}
	}
	// loopBack node
	nodeChunk.Write(IntToBinaryBE(p.Meta.NodeCount, 32))
	nodeChunk.Write(IntToBinaryBE(p.Meta.NodeCount, 32))

	// MetaData Chunk
	metaDataChunk := &bytes.Buffer{}
	p.Meta.TotalSize = nodeChunk.Len() + p.dataChunk.Len()
	metaData, err := json.Marshal(p.Meta)
	if err != nil {
		return err
	}
	metaDataChunk.Write(IntToBinaryBE(len(metaData), 32))
	metaDataChunk.Write(metaData)

	// Result
	if _, err := metaDataChunk.WriteTo(w); err != nil {
		return err
	}
	if _, err := nodeChunk.WriteTo(w); err != nil {
		return err
	}
	if _, err := p.dataChunk.WriteTo(w); err != nil {
		return err
	}

	return nil
}

// insert 插入数据
func (p *Writer) insert(ipNet *net.IPNet, values []string) error {
	mask, _ := ipNet.Mask.Size()
	node, index, ok := p.Nodes(ipNet.IP, mask)
	if !ok {
		log.Printf("load cidr failed cidr(%s) data(%s) node(%d) index(%d) preview data(%s)\n", ipNet, values, node, index, p.resolve(-p.node[node][index]))
		return errors.ErrInvalidCIDR
	}
	if p.node[node][index] > 0 {
		log.Printf("cidr conflict %s %s\n", ipNet, values)
		return errors.ErrCIDRConflict
	}
	offset := p.Fields(values)
	p.node[node][index] = -offset
	return nil
}

// resolve 解析数据
func (p *Writer) resolve(offset int) string {
	offset -= 8
	data := p.dataChunk.Bytes()
	size := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	if (offset + 2 + size) > len(data) {
		return ""
	}
	return string(data[offset+2 : offset+2+size])
}

// Nodes 获取CIDR地址所在节点和index
// 将补全Node中间链路，如果中间链路已经有数据，将无法写入新数据
func (p *Writer) Nodes(ip net.IP, mask int) (node, index int, ok bool) {
	// 如果传入的是IPv4，子网掩码增加96位( len(IPv6)-len(IPv4) )
	// 统一扩展为IPv6的子网掩码进行处理
	maxMask := mask - 1
	if ip.To4() != nil {
		if maxMask < 32 {
			maxMask += 96
		}
		if len(ip) == net.IPv4len {
			ip = ip.To16()
		}
	}
	for i := 0; i < maxMask; i++ {
		index = ((0xFF & int(ip[i>>3])) >> uint(7-(i%8))) & 1
		if p.node[node][index] == 0 {
			p.node = append(p.node, [2]int{})
			p.node[node][index] = len(p.node) - 1
		}
		if p.node[node][index] < 0 {
			return node, index, false
		}
		node = p.node[node][index]
	}
	return node, ((0xFF & int(ip[maxMask>>3])) >> uint(7-(maxMask%8))) & 1, true
}

// Fields 保存数据并返回数据的偏移量
// 相同的数据仅保存一份
// 数据格式 2 byte length + n byte data
func (p *Writer) Fields(fields []string) int {
	data := strings.Join(fields, FieldsSep)
	if _, ok := p.dataHash[data]; !ok {
		_data := []byte(data)
		// +8 是由于 loopBack node 占用了8byte
		p.dataHash[data] = p.dataChunk.Len() + 8
		p.dataChunk.Write(IntToBinaryBE(len(_data), 16))
		p.dataChunk.Write(_data)
	}
	return p.dataHash[data]
}

// IntToBinaryBE 将int转换为 binary big endian
func IntToBinaryBE(num, length int) []byte {
	switch length {
	case 16:
		_num := uint16(num)
		return []byte{byte(_num >> 8), byte(_num)}
	case 32:
		_num := uint32(num)
		return []byte{byte(_num >> 24), byte(_num >> 16), byte(_num >> 8), byte(_num)}
	default:
		return []byte{}
	}
}
