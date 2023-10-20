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

package ipdb

import (
	"net"

	"github.com/sjzar/ips/format/ipdb/sdk"
	"github.com/sjzar/ips/ipnet"
	"github.com/sjzar/ips/pkg/model"
)

const (
	DBFormat = "ipdb"
	DBExt    = ".ipdb"
)

// Reader is a structure that provides functionalities to read from IPDB IP database.
type Reader struct {
	meta *model.Meta // Metadata of the IP database
	db   *sdk.City   // Database reader instance
}

// NewReader initializes a new instance of Reader.
func NewReader(file string) (*Reader, error) {
	city, err := sdk.NewCity(file)
	if err != nil {
		return nil, err
	}

	meta := &model.Meta{
		MetaVersion: model.MetaVersion,
		Format:      DBFormat,
		Fields:      city.Fields(),
	}
	meta.AddCommonFieldAlias(CommonFieldsAlias)

	if city.IsIPv4() {
		meta.IPVersion |= model.IPv4
	}
	if city.IsIPv6() {
		meta.IPVersion |= model.IPv6
	}

	return &Reader{
		meta: meta,
		db:   city,
	}, nil
}

// Find retrieves IP information based on the given IP address.
func (r *Reader) Find(ip net.IP) (*model.IPInfo, error) {
	data, ipNet, err := r.db.FindMap(ip.String(), "CN")
	if err != nil {
		return nil, err
	}

	ret := &model.IPInfo{
		IP:     ip,
		IPNet:  ipnet.NewRange(ipNet),
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
