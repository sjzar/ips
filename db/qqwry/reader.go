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

package qqwry

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
)

const (
	// RedirectMode1 重定向模式1
	// 表示国家记录和地区记录都被重定向
	RedirectMode1 = 0x01

	// RedirectMode2 重定向模式2
	// 表示国家记录或地区记录被重定向
	RedirectMode2 = 0x02
)

/* qqwry.dat Format (Little Endian + GBK Encoding)
+--------------------------------+--------------------------------+
|       Start Index (4byte)      |       End Index (4byte)        |
+--------------------------------+--------------------------------+
|                            Data Chunk                           |
+--------------------------------+--------------------------------+
|                            Index Chunk                          |
+--------------------------------+--------------------------------+

Data Chunk
+--------------------------------+--------------------------------+
|                           End IP (4byte)                        |
+--------------------------------+--------------------------------+
|         Country (n byte)       |     End Flag 0x00 (1byte)      |
+--------------------------------+--------------------------------+
|          Area (n byte)         |     End Flag 0x00 (1byte)      |
+--------------------------------+--------------------------------+

Redirect Mode1: redirect country AND area
+--------------------------------+--------------------------------+
|                           End IP (4byte)                        |
+--------------------------------+--------------------------------+
|   Redirect Mode1 0x01 (1byte)  |      Data Offset (3byte)       |
+--------------------------------+--------------------------------+

Redirect Mode2: redirect country OR area
+--------------------------------+--------------------------------+
|                           End IP (4byte)                        |
+--------------------------------+--------------------------------+
|   Redirect Mode2 0x02 (1byte)  |      Data Offset (3byte)       |
+--------------------------------+--------------------------------+
|          Area (n byte)         |     End Flag 0x00 (1byte)      |
+--------------------------------+--------------------------------+

Index Chunk
+--------------------------------+--------------------------------+
|        Start IP (4byte)        |       Data Offset (3byte)      |
+--------------------------------+--------------------------------+

*/

// QQWry QQWry 数据库
type QQWry struct {

	// data IP库数据
	data []byte

	// start IP库数据开始位置
	start uint32

	// end IP库数据结束位置
	end uint32

	// gbkDecoder GBK解码器
	gbkDecoder *encoding.Decoder
}

func NewQQWry(filePath string) (*QQWry, error) {
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

	if len(data) < 8 {
		return nil, errors.ErrDatabaseIsInvalid
	}

	start := binary.LittleEndian.Uint32(data[:4])
	end := binary.LittleEndian.Uint32(data[4:])

	if uint32(len(data)) < end+7 {
		return nil, errors.ErrDatabaseIsInvalid
	}

	return &QQWry{
		data:       data,
		start:      start,
		end:        end,
		gbkDecoder: simplifiedchinese.GBK.NewDecoder(),
	}, nil
}

// Find 查找IP
func (q *QQWry) Find(ip net.IP) (*ipx.Range, map[string]string, error) {

	ip = ip.To4()
	if ip == nil {
		return nil, nil, errors.ErrIPVersionNotSupported
	}

	startIP, offset := q.findOffset(ipx.IPv4ToUint32(ip))
	if offset == 0 {
		return nil, nil, errors.ErrDatabaseIsInvalid
	}
	endIP := binary.LittleEndian.Uint32(q.data[offset : offset+4])
	country, area, err := q.parse(offset+4, 0)
	if err != nil {
		return nil, nil, err
	}

	return &ipx.Range{
			Start: ipx.Uint32ToIPv4(startIP).To4(),
			End:   ipx.Uint32ToIPv4(endIP).To4(),
		}, map[string]string{
			FieldCountry: country,
			FieldArea:    area,
		}, nil
}

// findOffset 查找IP对应的偏移量
func (q *QQWry) findOffset(ip uint32) (startIP uint32, offset uint32) {
	low := q.start
	high := q.end

	var mid, currentIP uint32
	var buf []byte

	for {
		mid = (high-low)/7/2*7 + low
		buf = q.data[mid : mid+7]
		currentIP = binary.LittleEndian.Uint32(buf[:4])

		if high-low == 7 {
			if binary.LittleEndian.Uint32(q.data[high:high+4]) <= ip {
				buf = q.data[high : high+7]
			}
			return binary.LittleEndian.Uint32(buf[:4]), Bytes3Uint32(buf[4:7])
		}

		if currentIP > ip {
			high = mid
		} else if currentIP < ip {
			low = mid
		} else {
			return ip, Bytes3Uint32(buf[4:7])
		}
	}
}

// parse 解析数据
func (q *QQWry) parse(offset uint32, depth int) (country, area string, err error) {
	if depth > 1 {
		return "", "", errors.ErrDatabaseIsInvalid
	}

	switch q.data[offset] {
	case RedirectMode1:
		// Redirect Mode1: redirect country AND area
		return q.parse(Bytes3Uint32(q.data[offset+1:offset+4]), depth+1)
	case RedirectMode2:
		// Redirect Mode2: redirect country OR area
		country, _, err = q.parseString(Bytes3Uint32(q.data[offset+1 : offset+4]))
		if err != nil {
			return "", "", err
		}
		offset += 4
	default:
		var length int
		country, length, err = q.parseString(offset)
		if err != nil {
			return "", "", err
		}
		// +1 跳过结束标志(0x00)
		offset += uint32(length) + 1
	}
	area, err = q.parseArea(offset, depth)
	if err != nil {
		return "", "", err
	}

	return country, area, nil
}

// parseArea 解析地区
func (q *QQWry) parseArea(offset uint32, depth int) (area string, err error) {
	if depth > 2 {
		return "", errors.ErrDatabaseIsInvalid
	}

	switch q.data[offset] {
	case RedirectMode1, RedirectMode2:
		return q.parseArea(Bytes3Uint32(q.data[offset+1:offset+4]), depth+1)
	}
	area, _, err = q.parseString(offset)
	if err != nil {
		return "", err
	}

	return area, nil
}

// parseString 解析字符串
func (q *QQWry) parseString(offset uint32) (string, int, error) {
	length := bytes.IndexByte(q.data[offset:], 0x00)
	if length == -1 {
		return "", 0, errors.ErrDatabaseIsInvalid
	}
	str, _ := q.gbkDecoder.String(string(q.data[offset : offset+uint32(length)]))
	return str, length, nil
}

// Bytes3Uint32 3字节转换为uint32
func Bytes3Uint32(b []byte) uint32 {
	_ = b[2]
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}
