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

	"github.com/sjzar/ips/format/mmdb/sdk"
	"github.com/sjzar/ips/pkg/model"
)

const (
	DBFormat = "mmdb"
	DBExt    = ".mmdb"
)

// Reader is a structure that provides functionalities to read from MMDB IP database.
type Reader struct {
	meta   *model.Meta  // Metadata of the IP database
	db     *sdk.Reader  // Database reader instance
	option ReaderOption // Configuration options for the reader.
}

// NewReader initializes and returns a new Reader for the specified MMDB file.
// It loads metadata and sets default options for the reader.
func NewReader(file string) (*Reader, error) {
	db, err := sdk.NewReader(file)
	if err != nil {
		return nil, err
	}

	meta := &model.Meta{
		MetaVersion: model.MetaVersion,
		Format:      DBFormat,
		IPVersion:   db.IPVersion,
		Fields:      db.Fields,
	}
	meta.AddCommonFieldAlias(CommonFieldsAlias)

	return &Reader{
		meta: meta,
		db:   db,
	}, nil
}

// Find retrieves IP information based on the given IP address.
func (r *Reader) Find(ip net.IP) (*model.IPInfo, error) {
	ipNet, data, err := r.db.Find(ip)
	if err != nil {
		return nil, err
	}

	ret := &model.IPInfo{
		IP:     ip,
		IPNet:  ipNet,
		Data:   data,
		Fields: r.meta.Fields,
	}
	ret.AddCommonFieldAlias(CommonFieldsAlias)

	return ret, nil
}

// Meta returns the meta-information of the IP database.
func (r *Reader) Meta() *model.Meta {
	return r.meta
}

// ReaderOption contains configuration options for the Reader.
type ReaderOption struct {
	DisableExtraData bool // If true, extra data (matched via GeoNameID) won't be used.
	UseFullField     bool // If true, all data will be combined into a single JSON string field.
}

// SetOption applies the provided option to the Reader's configuration.
func (r *Reader) SetOption(option interface{}) error {
	if opt, ok := option.(ReaderOption); ok {
		r.db.DisableExtraData = opt.DisableExtraData
		r.db.UseFullField = opt.UseFullField
		r.option = opt
	}
	return nil
}

// Close releases any resources used by the Reader and closes the MMDB database.
func (r *Reader) Close() error {
	return r.db.Close()
}
