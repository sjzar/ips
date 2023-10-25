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
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

// Copy From https://github.com/ipipdotnet/ipdb-go
// modify by shenjunzheng
// modify: Find / FindMap return ipNet

const IPv4 = 0x01
const IPv6 = 0x02

var (
	ErrFileSize = errors.New("IP Database file size error")
	ErrMetaData = errors.New("IP Database metadata error")
	ErrReadFull = errors.New("IP Database ReadFull error")

	ErrDatabaseError = errors.New("database error")

	ErrIPFormat = errors.New("query IP format error")

	ErrNoSupportLanguage = errors.New("language not support")
	ErrNoSupportIPv4     = errors.New("IPv4 not support")
	ErrNoSupportIPv6     = errors.New("IPv6 not support")

	ErrDataNotExists = errors.New("data is not exists")
)

type MetaData struct {
	Build     int64          `json:"build"`
	IPVersion uint16         `json:"ip_version"`
	Languages map[string]int `json:"languages"`
	NodeCount int            `json:"node_count"`
	TotalSize int            `json:"total_size"`
	Fields    []string       `json:"fields"`
}

type reader struct {
	fileSize  int
	nodeCount int
	v4offset  int

	meta MetaData
	data []byte

	refType map[string]string
}

func newReader(name string, obj interface{}) (*reader, error) {
	var err error
	var fileInfo os.FileInfo
	fileInfo, err = os.Stat(name)
	if err != nil {
		return nil, err
	}
	fileSize := int(fileInfo.Size())
	if fileSize < 4 {
		return nil, ErrFileSize
	}
	body, err := os.ReadFile(name)
	if err != nil {
		return nil, ErrReadFull
	}
	var meta MetaData
	metaLength := int(binary.BigEndian.Uint32(body[0:4]))
	if fileSize < (4 + metaLength) {
		return nil, ErrFileSize
	}
	if err := json.Unmarshal(body[4:4+metaLength], &meta); err != nil {
		return nil, err
	}
	if len(meta.Languages) == 0 || len(meta.Fields) == 0 {
		return nil, ErrMetaData
	}
	if fileSize != (4 + metaLength + meta.TotalSize) {
		return nil, ErrFileSize
	}

	var dm map[string]string
	if obj != nil {
		t := reflect.TypeOf(obj).Elem()
		dm = make(map[string]string, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			k := t.Field(i).Tag.Get("json")
			dm[k] = t.Field(i).Name
		}
	}

	db := &reader{
		fileSize:  fileSize,
		nodeCount: meta.NodeCount,

		meta:    meta,
		refType: dm,

		data: body[4+metaLength:],
	}

	if db.v4offset == 0 {
		node := 0
		for i := 0; i < 96 && node < db.nodeCount; i++ {
			if i >= 80 {
				node = db.readNode(node, 1)
			} else {
				node = db.readNode(node, 0)
			}
		}
		db.v4offset = node
	}

	return db, nil
}

func newIOReader(r io.Reader, obj interface{}) (*reader, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, ErrReadFull
	}

	fileSize := len(body)
	var meta MetaData
	metaLength := int(binary.BigEndian.Uint32(body[0:4]))
	if fileSize < (4 + metaLength) {
		return nil, ErrFileSize
	}
	if err := json.Unmarshal(body[4:4+metaLength], &meta); err != nil {
		return nil, err
	}
	if len(meta.Languages) == 0 || len(meta.Fields) == 0 {
		return nil, ErrMetaData
	}
	if fileSize != (4 + metaLength + meta.TotalSize) {
		return nil, ErrFileSize
	}

	var dm map[string]string
	if obj != nil {
		t := reflect.TypeOf(obj).Elem()
		dm = make(map[string]string, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			k := t.Field(i).Tag.Get("json")
			dm[k] = t.Field(i).Name
		}
	}

	db := &reader{
		fileSize:  fileSize,
		nodeCount: meta.NodeCount,

		meta:    meta,
		refType: dm,

		data: body[4+metaLength:],
	}

	if db.v4offset == 0 {
		node := 0
		for i := 0; i < 96 && node < db.nodeCount; i++ {
			if i >= 80 {
				node = db.readNode(node, 1)
			} else {
				node = db.readNode(node, 0)
			}
		}
		db.v4offset = node
	}

	return db, nil
}

func (db *reader) Find(addr, language string) ([]string, *net.IPNet, error) {
	return db.find1(addr, language)
}

func (db *reader) FindMap(addr, language string) (map[string]string, *net.IPNet, error) {

	data, ipNet, err := db.find1(addr, language)
	if err != nil {
		return nil, nil, err
	}
	info := make(map[string]string, len(db.meta.Fields))
	for k, v := range data {
		info[db.meta.Fields[k]] = v
	}

	return info, ipNet, nil
}

func (db *reader) find0(addr string) ([]byte, *net.IPNet, error) {

	var err error
	var node int
	var mask int
	var ipNet *net.IPNet
	ipv := net.ParseIP(addr)
	if ip := ipv.To4(); ip != nil {
		if !db.IsIPv4Support() {
			return nil, nil, ErrNoSupportIPv4
		}

		node, mask, err = db.search(ip, 32)
		cidrMask := net.CIDRMask(mask, len(ip)*8)
		ipNet = &net.IPNet{IP: ipv.Mask(cidrMask), Mask: cidrMask}
	} else if ip := ipv.To16(); ip != nil {
		if !db.IsIPv6Support() {
			return nil, nil, ErrNoSupportIPv6
		}

		node, mask, err = db.search(ip, 128)
		cidrMask := net.CIDRMask(mask, len(ip)*8)
		ipNet = &net.IPNet{IP: ipv.Mask(cidrMask), Mask: cidrMask}
	} else {
		return nil, nil, ErrIPFormat
	}

	if err != nil || node < 0 {
		return nil, nil, err
	}

	body, err := db.resolve(node)
	if err != nil {
		return nil, nil, err
	}

	return body, ipNet, nil
}

func (db *reader) find1(addr, language string) ([]string, *net.IPNet, error) {

	off, ok := db.meta.Languages[language]
	if !ok {
		return nil, nil, ErrNoSupportLanguage
	}

	body, ipNet, err := db.find0(addr)
	if err != nil {
		return nil, nil, err
	}

	str := (*string)(unsafe.Pointer(&body))
	tmp := strings.Split(*str, "\t")

	if (off + len(db.meta.Fields)) > len(tmp) {
		return nil, nil, ErrDatabaseError
	}

	return tmp[off : off+len(db.meta.Fields)], ipNet, nil
}

func (db *reader) search(ip net.IP, bitCount int) (int, int, error) {

	var node int

	if bitCount == 32 {
		node = db.v4offset
	} else {
		node = 0
	}

	var i = 0
	for ; i < bitCount; i++ {
		if node >= db.nodeCount {
			break
		}

		node = db.readNode(node, ((0xFF&int(ip[i>>3]))>>uint(7-(i%8)))&1)
	}

	if node > db.nodeCount {
		return node, i, nil
	}

	return -1, 0, ErrDataNotExists
}

func (db *reader) readNode(node, index int) int {
	off := node*8 + index*4
	return int(binary.BigEndian.Uint32(db.data[off : off+4]))
}

func (db *reader) resolve(node int) ([]byte, error) {
	resolved := node - db.nodeCount + db.nodeCount*8
	if resolved >= db.fileSize {
		return nil, ErrDatabaseError
	}

	size := int(binary.BigEndian.Uint16(db.data[resolved : resolved+2]))
	if (resolved + 2 + size) > len(db.data) {
		return nil, ErrDatabaseError
	}
	bytes := db.data[resolved+2 : resolved+2+size]

	return bytes, nil
}

func (db *reader) IsIPv4Support() bool {
	return (int(db.meta.IPVersion) & IPv4) == IPv4
}

func (db *reader) IsIPv6Support() bool {
	return (int(db.meta.IPVersion) & IPv6) == IPv6
}

func (db *reader) Build() time.Time {
	return time.Unix(db.meta.Build, 0).In(time.UTC)
}

func (db *reader) Languages() []string {
	ls := make([]string, 0, len(db.meta.Languages))
	for k := range db.meta.Languages {
		ls = append(ls, k)
	}
	return ls
}
