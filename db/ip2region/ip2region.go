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

package ip2region

import (
	"net"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

type Database struct {
	meta model.Meta
	db   *IP2Region
}

// New 初始化 IP2Region 数据库实例
func New(file string) (*Database, error) {

	db, err := NewIP2Region(file)
	if err != nil {
		return nil, err
	}

	meta := model.Meta{
		Fields:    FullFields,
		IPVersion: model.IPv4,
	}

	return &Database{
		meta: meta,
		db:   db,
	}, nil
}

// Find 查询 IP 对应的网段和结果
func (d *Database) Find(ip net.IP) (*ipx.Range, map[string]string, error) {
	ipr, data, err := d.db.Find(ip)
	if err != nil {
		return nil, nil, err
	}
	return ipr, FieldsFormat(data), err
}

// Meta 返回元数据
func (d *Database) Meta() model.Meta {
	return d.meta
}

// Close 关闭数据库实例
func (d *Database) Close() error {
	return nil
}
