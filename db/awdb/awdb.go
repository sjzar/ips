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

package awdb

import (
	"net"

	"github.com/dilfish/awdb-golang/awdb-golang"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

// Database AWDB 数据库
type Database struct {
	meta model.Meta
	db   *awdb.Reader
}

// New 初始化 AWDB 数据库实例
func New(file string) (*Database, error) {
	db, err := awdb.Open(file)
	if err != nil {
		return nil, err
	}

	meta := model.Meta{
		Fields: FullFields,
	}
	if db.Metadata.IPVersion == 4 {
		meta.IPVersion |= model.IPv4
	}
	if db.Metadata.IPVersion == 6 {
		meta.IPVersion |= model.IPv6
	}

	return &Database{
		meta: meta,
		db:   db,
	}, nil
}

// Meta 返回元数据
func (d *Database) Meta() model.Meta {
	return d.meta
}

// Find 查询 IP 对应的网段和结果
func (d *Database) Find(ip net.IP) (*ipx.Range, map[string]string, error) {
	var record interface{}
	ipNet, _, err := d.db.LookupNetwork(ip, &record)
	if err != nil {
		return nil, nil, err
	}

	data := make(map[string]string)
	for k, v := range record.(map[string]interface{}) {
		data[k] = string(v.([]byte))
	}

	return ipx.NewRange(ipNet), model.FieldsFormat(CommonFieldsMap, data), nil
}

// Close 关闭数据库实例
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}
