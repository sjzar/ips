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

package awdb

import (
	"net"

	"github.com/dilfish/awdb-golang/awdb-golang"

	"github.com/sjzar/ips/ipnet"
	"github.com/sjzar/ips/pkg/model"
)

const (
	DBFormat = "awdb"
	DBExt    = ".awdb"
)

// Reader is a structure that provides functionalities to read from AWDB IP database.
type Reader struct {
	meta *model.Meta  // Metadata of the IP database
	db   *awdb.Reader // Database reader instance
}

// NewReader initializes a new instance of Reader.
func NewReader(file string) (*Reader, error) {
	db, err := awdb.Open(file)
	if err != nil {
		return nil, err
	}

	meta := &model.Meta{
		MetaVersion: model.MetaVersion,
		Format:      DBFormat,
		Fields:      FullFields,
	}
	meta.AddCommonFieldAlias(CommonFieldsAlias)

	if db.Metadata.IPVersion == 4 {
		meta.IPVersion |= model.IPv4
	}
	if db.Metadata.IPVersion == 6 {
		meta.IPVersion |= model.IPv6
	}

	return &Reader{
		meta: meta,
		db:   db,
	}, nil
}

// Meta returns the meta-information of the IP database.
func (r *Reader) Meta() *model.Meta {
	return r.meta
}

// Find retrieves IP information based on the given IP address.
func (r *Reader) Find(ip net.IP) (*model.IPInfo, error) {

	var record interface{}
	ipNet, _, err := r.db.LookupNetwork(ip, &record)
	if err != nil {
		return nil, err
	}

	ret := &model.IPInfo{
		IP:     ip,
		IPNet:  ipnet.NewRange(ipNet),
		Fields: r.meta.Fields,
	}

	ret.Data = make(map[string]string)
	for k, _v := range record.(map[string]interface{}) {
		ret.Data[k] = string(_v.([]byte))
	}
	ret.AddCommonFieldAlias(CommonFieldsAlias)

	return ret, nil
}

// SetOption configures the Reader with the provided option.
func (r *Reader) SetOption(option interface{}) error {
	return nil
}

// Close closes the IP database.
func (r *Reader) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
