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
	"net"

	"github.com/sjzar/ips/db"
	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

type DBReader struct {

	// meta 元数据
	meta model.Meta

	// database IP库
	db db.Database

	// selector 字段选择工具
	selector *FieldSelector
}

// NewDBReader 初始化 IP 库读取器
func NewDBReader(db db.Database, selector *FieldSelector) *DBReader {
	meta := model.Meta{
		IPVersion: db.Meta().IPVersion,
		Fields:    selector.Fields(),
	}

	return &DBReader{
		meta:     meta,
		db:       db,
		selector: selector,
	}
}

// Meta 返回 meta 信息
func (s *DBReader) Meta() model.Meta {
	return s.meta
}

// Find 查询IP所在网段和对应数据
func (s *DBReader) Find(ip net.IP) (*ipx.Range, []string, error) {
	ipRange, data, err := s.db.Find(ip)
	if err != nil {
		return nil, nil, err
	}

	return ipRange, s.selector.Select(data), nil
}

// Close 关闭
func (s *DBReader) Close() error {
	return s.db.Close()
}
