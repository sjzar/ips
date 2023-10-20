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

package mmdb

import (
	"net"

	"github.com/oschwald/maxminddb-golang"

	"github.com/sjzar/ips/ipnet"
	"github.com/sjzar/ips/pkg/model"
)

const (
	DBFormat = "mmdb"
	DBExt    = ".mmdb"
)

// Reader is a structure that provides functionalities to read from MMDB IP database.
type Reader struct {
	meta   *model.Meta       // Metadata of the IP database
	db     *maxminddb.Reader // Database reader instance
	dbType databaseType      // Database type
}

// NewReader initializes a new instance of Reader.
func NewReader(file string) (*Reader, error) {

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

	meta := &model.Meta{
		MetaVersion: model.MetaVersion,
		Format:      DBFormat,
		Fields:      fullFields,
	}
	meta.AddCommonFieldAlias(CommonFieldsAlias)

	switch db.Metadata.IPVersion {
	case 4:
		meta.IPVersion |= model.IPv4
	case 6:
		meta.IPVersion |= model.IPv4
		meta.IPVersion |= model.IPv6
	}

	return &Reader{
		meta:   meta,
		db:     db,
		dbType: dbType,
	}, nil
}

// Find retrieves IP information based on the given IP address.
func (d *Reader) Find(ip net.IP) (*model.IPInfo, error) {

	data := make(map[string]string)
	var ipNet *net.IPNet
	var err error
	switch {
	case d.dbType&isCity > 0:
		var city City
		ipNet, _, err = d.db.LookupNetwork(ip, &city)
		if err != nil {
			return nil, err
		}
		data = city.Format()
	case d.dbType&isASN > 0:
		var asn ASN
		ipNet, _, err = d.db.LookupNetwork(ip, &asn)
		if err != nil {
			return nil, err
		}
		data = asn.Format()
	}

	ret := &model.IPInfo{
		IP:     ip,
		IPNet:  ipnet.NewRange(ipNet),
		Data:   data,
		Fields: d.meta.Fields,
	}
	ret.AddCommonFieldAlias(CommonFieldsAlias)

	return ret, nil
}

// Meta returns the meta-information of the IP database.
func (d *Reader) Meta() *model.Meta {
	return d.meta
}

// SetOption configures the Reader with the provided option.
func (d *Reader) SetOption(option interface{}) error {
	return nil
}

// Close closes the IP database.
func (d *Reader) Close() error {
	return d.db.Close()
}
