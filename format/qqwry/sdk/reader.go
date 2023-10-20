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

package sdk

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/sjzar/ips/ipnet"
	"github.com/sjzar/ips/pkg/errors"
)

const (
	// RedirectMode1 重定向模式1
	// 表示国家记录和地区记录都被重定向
	RedirectMode1 = 0x01

	// RedirectMode2 重定向模式2
	// 表示国家记录或地区记录被重定向
	RedirectMode2 = 0x02
)

// Reader represents the QQWry database reader.
type Reader struct {
	data       []byte            // IP database data
	start      uint32            // Start position of IP database data
	end        uint32            // End position of IP database data
	gbkDecoder *encoding.Decoder // Decoder for GBK encoding
}

// NewReader initializes a new QQWry instance given the file path.
func NewReader(filePath string) (*Reader, error) {
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
		return nil, errors.ErrInvalidDatabase
	}

	start := binary.LittleEndian.Uint32(data[:4])
	end := binary.LittleEndian.Uint32(data[4:])

	if uint32(len(data)) < end+7 {
		return nil, errors.ErrInvalidDatabase
	}

	return &Reader{
		data:       data,
		start:      start,
		end:        end,
		gbkDecoder: simplifiedchinese.GBK.NewDecoder(),
	}, nil
}

// Find locates the IP in the QQWry database and returns its range, country, and area.
func (q *Reader) Find(ip net.IP) (*ipnet.Range, string, string, error) {

	ip = ip.To4()
	if ip == nil {
		return nil, "", "", errors.ErrUnsupportedIPVersion
	}

	startIP, offset := q.findOffset(ipnet.IPv4ToUint32(ip))
	if offset == 0 {
		return nil, "", "", errors.ErrInvalidDatabase
	}
	endIP := binary.LittleEndian.Uint32(q.data[offset : offset+4])
	country, area, err := q.parse(offset+4, 0)
	if err != nil {
		return nil, "", "", err
	}

	return &ipnet.Range{
		Start: ipnet.Uint32ToIPv4(startIP).To4(),
		End:   ipnet.Uint32ToIPv4(endIP).To4(),
	}, country, area, nil
}

// findOffset determines the offset for the given IP in the QQWry database.
func (q *Reader) findOffset(ip uint32) (startIP uint32, offset uint32) {
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

// parse extracts the country and area data for the given offset in the QQWry database.
func (q *Reader) parse(offset uint32, depth int) (country, area string, err error) {
	if depth > 1 {
		return "", "", errors.ErrInvalidDatabase
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

// parseArea retrieves the area data for the given offset in the QQWry database.
func (q *Reader) parseArea(offset uint32, depth int) (area string, err error) {
	if depth > 2 {
		return "", errors.ErrInvalidDatabase
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

// parseString decodes and retrieves the string data for the given offset in the QQWry database.
func (q *Reader) parseString(offset uint32) (string, int, error) {
	length := bytes.IndexByte(q.data[offset:], 0x00)
	if length == -1 {
		return "", 0, errors.ErrInvalidDatabase
	}
	str, _ := q.gbkDecoder.String(string(q.data[offset : offset+uint32(length)]))
	return str, length, nil
}

// Bytes3Uint32 converts a 3-byte slice to a uint32 value.
func Bytes3Uint32(b []byte) uint32 {
	_ = b[2]
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}
