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

// Copy From https://github.com/lionsoul2014/ip2region/blob/master/binding/golang/xdb/searcher.go
// Author: lionsoul2014

package ip2region

import (
	"encoding/binary"
	"io"
	"net"
	"os"
	"strings"

	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
)

/* IP2Region Format (Little Endian)
+--------------------------------+--------------------------------+
|                      Header Chunk (256byte)                     |
+--------------------------------+--------------------------------+
|                  Vector Index Chunk (256*256byte)               |
+--------------------------------+--------------------------------+
|                           Data Chunk                            |
+--------------------------------+--------------------------------+
|                           Index Chunk                           |
+--------------------------------+--------------------------------+

Header Chunk (256byte)
+--------------------------------+--------------------------------+
| Version (2byte) | Cache Policy (2byte) |   Build Time (4byte)   |
+--------------------------------+--------------------------------+
|      Start Index (4byte)       |       End Index (4byte)        |
+--------------------------------+--------------------------------+
|                           Empty Data                            |
+--------------------------------+--------------------------------+

Vector Index Chunk (256*256byte)
+--------------------------------+--------------------------------+
|   Index Start Offset (4byte)   |    Index End Offset (4byte)    |
+--------------------------------+--------------------------------+

Index Chunk
+--------------------------------+--------------------------------+
|        Start IP (4byte)        |         End IP (4byte)         |
+--------------------------------+--------------------------------+
|       Data Length (2byte)      |       Data Offset (4byte)      |
+--------------------------------+--------------------------------+

Document: https://mp.weixin.qq.com/s/ndjzu0BgaeBmDOCw5aqHUg
*/

const (
	// HeaderInfoLength 长度为 256 字节
	HeaderInfoLength = 256

	// VectorIndexCols Vector 索引列数
	VectorIndexCols = 256

	// VectorIndexSize Vector 索引长度
	VectorIndexSize = 8

	// IndexLen 索引长度
	IndexLen = 14

	// FieldSpe 字段分隔符
	FieldSpe = "|"
)

type IP2Region struct {
	data []byte
}

func NewIP2Region(filePath string) (*IP2Region, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) < 256 {
		return nil, errors.ErrDatabaseIsInvalid
	}

	end := binary.LittleEndian.Uint32(data[12:16])
	if end+IndexLen > uint32(len(data)) {
		return nil, errors.ErrDatabaseIsInvalid
	}

	return &IP2Region{
		data: data,
	}, nil
}

// Find 查找IP
func (i *IP2Region) Find(ip net.IP) (*ipx.Range, map[string]string, error) {

	ip = ip.To4()
	if ip == nil {
		return nil, nil, errors.ErrIPVersionNotSupported
	}

	startIP, endIP, length, offset := i.findOffset(ipx.IPv4ToUint32(ip))
	data := make(map[string]string)
	if endIP == 0 || offset == 0 {
		return nil, nil, errors.ErrDatabaseIsInvalid
	}

	if length != 0 {
		str := string(i.data[offset : offset+length])
		split := strings.Split(str, FieldSpe)
		for i := range split {
			if i > len(FullFields) {
				break
			}
			if split[i] != "0" {
				data[FullFields[i]] = split[i]
			} else {
				data[FullFields[i]] = ""
			}
		}
	}

	return &ipx.Range{
		Start: ipx.Uint32ToIPv4(startIP).To4(),
		End:   ipx.Uint32ToIPv4(endIP).To4(),
	}, data, nil
}

// findOffset 查找IP对应的偏移量
func (i *IP2Region) findOffset(ip uint32) (startIP, endIP uint32, length, offset uint32) {

	// locate the segment index block based on the vector index
	var il0 = (ip >> 24) & 0xFF
	var il1 = (ip >> 16) & 0xFF
	var idx = il0*VectorIndexCols*VectorIndexSize + il1*VectorIndexSize
	var sPtr, ePtr = uint32(0), uint32(0)
	sPtr = binary.LittleEndian.Uint32(i.data[HeaderInfoLength+idx:])
	ePtr = binary.LittleEndian.Uint32(i.data[HeaderInfoLength+idx+4:])

	var l, h = 0, int((ePtr - sPtr) / IndexLen)
	for l <= h {
		m := (l + h) >> 1
		p := sPtr + uint32(m*IndexLen)
		buff := i.data[p : p+IndexLen]

		// decode the data step by step to reduce the unnecessary operations
		startIP = binary.LittleEndian.Uint32(buff)
		if ip < startIP {
			h = m - 1
		} else {
			endIP = binary.LittleEndian.Uint32(buff[4:])
			if ip > endIP {
				l = m + 1
			} else {
				length = uint32(binary.LittleEndian.Uint16(buff[8:]))
				offset = binary.LittleEndian.Uint32(buff[10:])
				break
			}
		}
	}

	return
}
