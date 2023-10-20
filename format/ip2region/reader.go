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

package ip2region

import (
	"net"

	"github.com/sjzar/ips/format/ip2region/sdk"
	"github.com/sjzar/ips/pkg/model"
)

const (
	DBFormat = "ip2region"
	DBExt    = ".xdb"
)

// Reader is a structure that provides functionalities to read from IP2Region IP database.
type Reader struct {
	meta *model.Meta // Metadata of the IP database
	db   *sdk.Reader // Database reader instance
}

// NewReader initializes a new instance of Reader.
func NewReader(file string) (*Reader, error) {

	db, err := sdk.NewReader(file)
	if err != nil {
		return nil, err
	}

	meta := &model.Meta{
		MetaVersion: model.MetaVersion,
		Format:      DBFormat,
		IPVersion:   model.IPv4,
		Fields:      FullFields,
	}
	meta.AddCommonFieldAlias(CommonFieldsAlias)

	return &Reader{
		meta: meta,
		db:   db,
	}, nil
}

// Find retrieves IP information based on the given IP address.
func (r *Reader) Find(ip net.IP) (*model.IPInfo, error) {
	ipr, values, err := r.db.Find(ip)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string, len(FullFields))
	for i := range values {
		if i > len(FullFields) {
			break
		}
		if values[i] != "0" {
			data[FullFields[i]] = values[i]
		} else {
			data[FullFields[i]] = ""
		}
	}

	ret := &model.IPInfo{
		IP:     ip,
		IPNet:  ipr,
		Fields: r.meta.Fields,
		Data:   data,
	}
	ret.AddCommonFieldAlias(CommonFieldsAlias)

	return ret, nil
}

// Meta returns the meta-information of the IP database.
func (r *Reader) Meta() *model.Meta {
	return r.meta
}

// SetOption configures the Reader with the provided option.
func (r *Reader) SetOption(option interface{}) error {
	return nil
}

// Close closes the IP database.
func (r *Reader) Close() error {
	return nil
}
