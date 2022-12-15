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

package mmdb

import (
	"net"

	"github.com/oschwald/maxminddb-golang"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

type Database struct {
	meta   model.Meta
	db     *maxminddb.Reader
	dbType databaseType
}

// New 初始化 mmdb 数据库实例
// 发现 GeoLite2-City 有数据不全的情况
func New(file string) (*Database, error) {

	db, err := maxminddb.Open(file)
	if err != nil {
		return nil, err
	}

	supportDefaultLang := false
	supportEnglish := false
	for _, lang := range db.Metadata.Languages {
		if lang == Lang {
			supportDefaultLang = true
			break
		}
		if lang == "en" {
			supportEnglish = true
		}
	}
	if !supportDefaultLang && supportEnglish {
		Lang = "en"
	}
	dbType := getDBType(db.Metadata.DatabaseType)
	fullFields := CityFullFields
	if dbType&isASN > 0 {
		fullFields = ASNFullFields
	}

	meta := model.Meta{
		Fields: fullFields,
	}

	switch db.Metadata.IPVersion {
	case 4:
		meta.IPVersion |= model.IPv4
	case 6:
		meta.IPVersion |= model.IPv4
		meta.IPVersion |= model.IPv6
	}

	return &Database{
		meta:   meta,
		db:     db,
		dbType: dbType,
	}, nil
}

// Find 查询 IP 对应的网段和结果
func (d *Database) Find(ip net.IP) (*ipx.Range, map[string]string, error) {

	dataMap := make(map[string]string)
	var ipNet *net.IPNet
	var err error
	switch {
	case d.dbType&isCity > 0:
		var data City
		ipNet, err = d.find(ip, &data)
		if err != nil {
			return nil, nil, err
		}
		dataMap = data.Format()
	case d.dbType&isASN > 0:
		var data ASN
		ipNet, err = d.find(ip, &data)
		if err != nil {
			return nil, nil, err
		}
		dataMap = data.Format()
	}

	return ipx.NewRange(ipNet), model.FieldsFormat(CommonFieldsMap, dataMap), nil
}

func (d *Database) find(ip net.IP, data interface{}) (*net.IPNet, error) {
	ipNet, _, err := d.db.LookupNetwork(ip, &data)
	if err != nil {
		return nil, err
	}
	return ipNet, nil
}

// Meta 返回元数据
func (d *Database) Meta() model.Meta {
	return d.meta
}

// Close 关闭数据库实例
func (d *Database) Close() error {
	return d.db.Close()
}
