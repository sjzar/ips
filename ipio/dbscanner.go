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

package ipio

import (
	"log"
	"net"
	"reflect"

	"github.com/sjzar/ips/db"
	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
	"github.com/sjzar/ips/rewriter"
)

// DBScanner IP库扫描读取工具
type DBScanner struct {

	// meta 元数据
	meta model.Meta

	// database IP库
	db db.Database

	// selector 字段选择工具
	selector *FieldSelector

	// rw 数据改写
	rw rewriter.Rewriter

	marker net.IP
	ipr    *ipx.Range
	values []string
	done   bool
	err    error
}

// NewDBScanner 初始化 IP 库扫描读取器
func NewDBScanner(db db.Database, selector *FieldSelector, rw rewriter.Rewriter) *DBScanner {
	meta := model.Meta{
		IPVersion: db.Meta().IPVersion,
		Fields:    selector.Fields(),
	}

	return &DBScanner{
		meta:     meta,
		db:       db,
		selector: selector,
		rw:       rw,
	}
}

// Meta 返回 meta 信息
func (s *DBScanner) Meta() model.Meta {
	return s.meta
}

// Find 查询IP所在网段和对应数据
func (s *DBScanner) Find(ip net.IP) (*ipx.Range, []string, error) {
	ipRange, data, err := s.db.Find(ip)
	if err != nil {
		return nil, nil, err
	}

	// Data Rewrite
	if s.rw != nil {
		_, ipRange, data, err = s.rw.Rewrite(ip, ipRange, data)
		if err != nil {
			return nil, nil, err
		}
	}

	return ipRange, s.selector.Select(data), nil
}

// Close 关闭
func (s *DBScanner) Close() error {
	return s.db.Close()
}

// Init 初始化
func (s *DBScanner) Init(meta model.Meta) error {
	s.meta = model.Meta{
		IPVersion: s.db.Meta().IPVersion,
		Fields:    s.selector.Fields(),
	}

	if meta.IPVersion != 0 {
		if meta.IPVersion&s.meta.IPVersion != meta.IPVersion {
			return errors.ErrIPVersionNotSupported
		}
		s.meta.IPVersion = meta.IPVersion
	}

	s.marker, s.done = nil, false
	return nil
}

// Scan 扫描IP库
func (s *DBScanner) Scan() bool {
	if s.done {
		return false
	}
	if s.marker == nil {
		if s.meta.IsIPv6Support() {
			s.marker = make(net.IP, net.IPv6len)
		} else {
			s.marker = net.IPv4(0, 0, 0, 0)
		}
	}
	s.ipr, s.values = nil, nil
	for {
		_ipr, _values, err := s.Find(s.marker)
		if err != nil {
			log.Println("DBScanner Find() failed ", s.marker, _ipr, _values, err)
			s.err = err
			return false
		}

		if s.ipr == nil {
			s.ipr, s.values = _ipr, _values
		} else {
			if !reflect.DeepEqual(s.values, _values) {
				break
			}
			if ok := s.ipr.Join(_ipr); !ok {
				break
			}
		}

		s.marker = ipx.NextIP(_ipr.End)
		if ipx.IsLastIP(_ipr.End, s.meta.IsIPv6Support()) {
			s.done = true
			break
		}
	}
	return true
}

func (s *DBScanner) Result() (*ipx.Range, []string) {
	return s.ipr, s.values
}

func (s *DBScanner) Err() error {
	return s.err
}
