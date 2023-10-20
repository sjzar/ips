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

package ipdb

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"strings"
	"time"

	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

const (
	FieldsSep = "\t"
)

// Writer provides functionalities to write IP data into IPDB format.
type Writer struct {
	meta      *model.Meta    // Metadata for the IP database
	ipdbMeta  Meta           // Metadata for IPDB format
	node      [][2]int       // Node data for IPDB format
	dataHash  map[string]int // Data hash for IPDB format
	dataChunk *bytes.Buffer  // Data chunk buffer for IPDB format
}

// NewWriter initializes a new Writer instance for writing IP data in IPDB format.
func NewWriter(meta *model.Meta) (*Writer, error) {
	return &Writer{
		meta: meta,
		ipdbMeta: Meta{
			Build:     int(time.Now().Unix()),
			IPVersion: meta.IPVersion,
			Languages: map[string]int{"CN": 0},
			Fields:    model.ConvertToDBFields(meta.Fields, meta.FieldAlias, CommonFieldsAlias),
		},
		node:      [][2]int{{}},
		dataChunk: &bytes.Buffer{},
		dataHash:  make(map[string]int),
	}, nil
}

// WriterOption provides options for the Writer.
type WriterOption struct {
	// Languages specifies multiple languages for the IPDB format.
	Languages map[string]int
}

// SetOption sets the provided options to the Writer.
func (w *Writer) SetOption(option interface{}) error {
	if opt, ok := option.(WriterOption); ok {
		if len(opt.Languages) > 0 {
			w.ipdbMeta.Languages = opt.Languages
		}
		return nil
	}

	return nil
}

// Insert adds the given IP information into the writer.
func (w *Writer) Insert(info *model.IPInfo) error {
	values := info.Values()
	if len(values) != len(w.ipdbMeta.Fields) {
		return errors.ErrMismatchedFieldsLength
	}

	for _, ipNet := range info.IPNet.IPNets() {
		if err := w.insert(ipNet, values); err != nil {
			return err
		}
	}

	return nil
}

// WriteTo writes the IP data into the provided writer in IPDB format.
func (w *Writer) WriteTo(iw io.Writer) (int64, error) {

	// Node Chunk
	nodeChunk := &bytes.Buffer{}
	w.ipdbMeta.NodeCount = len(w.node)
	for i := 0; i < w.ipdbMeta.NodeCount; i++ {
		for j := 0; j < 2; j++ {
			// 小于0: 数据记录，设置为NodeLength+Data偏移量
			// 等于0: 空值，设置为NodeLength
			// 大于0: Node跳转记录，不做调整
			if w.node[i][j] <= 0 {
				w.node[i][j] = w.ipdbMeta.NodeCount - w.node[i][j]
			}
			nodeChunk.Write(IntToBinaryBE(w.node[i][j], 32))
		}
	}
	// loopBack node
	nodeChunk.Write(IntToBinaryBE(w.ipdbMeta.NodeCount, 32))
	nodeChunk.Write(IntToBinaryBE(w.ipdbMeta.NodeCount, 32))

	// MetaData Chunk
	metaDataChunk := &bytes.Buffer{}
	w.ipdbMeta.TotalSize = nodeChunk.Len() + w.dataChunk.Len()
	metaData, err := json.Marshal(w.ipdbMeta)
	if err != nil {
		return 0, err
	}
	metaDataChunk.Write(IntToBinaryBE(len(metaData), 32))
	metaDataChunk.Write(metaData)

	// Result
	if _, err := metaDataChunk.WriteTo(iw); err != nil {
		return 0, err
	}
	if _, err := nodeChunk.WriteTo(iw); err != nil {
		return 0, err
	}
	if _, err := w.dataChunk.WriteTo(iw); err != nil {
		return 0, err
	}

	return 0, nil
}

// insert 插入数据
func (w *Writer) insert(ipNet *net.IPNet, values []string) error {
	mask, _ := ipNet.Mask.Size()
	node, index, ok := w.Nodes(ipNet.IP, mask)
	if !ok {
		// log.Printf("load cidr failed cidr(%s) data(%s) node(%d) index(%d) preview data(%s)\n", ipNet, values, node, index, w.resolve(-w.node[node][index]))
		return errors.ErrInvalidCIDR
	}
	if w.node[node][index] > 0 {
		// log.Printf("cidr conflict %s %s\n", ipNet, values)
		return errors.ErrCIDROverlap
	}
	offset := w.Fields(values)
	w.node[node][index] = -offset
	return nil
}

// Resolve 解析数据
func (w *Writer) Resolve(offset int) string {
	offset -= 8
	data := w.dataChunk.Bytes()
	size := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	if (offset + 2 + size) > len(data) {
		return ""
	}
	return string(data[offset+2 : offset+2+size])
}

// Nodes 获取CIDR地址所在节点和index
// 将补全Node中间链路，如果中间链路已经有数据，将无法写入新数据
func (w *Writer) Nodes(ip net.IP, mask int) (node, index int, ok bool) {
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
		if w.node[node][index] == 0 {
			w.node = append(w.node, [2]int{})
			w.node[node][index] = len(w.node) - 1
		}
		if w.node[node][index] < 0 {
			return node, index, false
		}
		node = w.node[node][index]
	}
	return node, ((0xFF & int(ip[maxMask>>3])) >> uint(7-(maxMask%8))) & 1, true
}

// Fields 保存数据并返回数据的偏移量
// 相同的数据仅保存一份
// 数据格式 2 byte length + n byte data
func (w *Writer) Fields(fields []string) int {
	data := strings.Join(fields, FieldsSep)
	if _, ok := w.dataHash[data]; !ok {
		_data := []byte(data)
		// +8 是由于 loopBack node 占用了8byte
		w.dataHash[data] = w.dataChunk.Len() + 8
		w.dataChunk.Write(IntToBinaryBE(len(_data), 16))
		w.dataChunk.Write(_data)
	}
	return w.dataHash[data]
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
