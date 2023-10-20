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

package ipio

import (
	"net"

	"github.com/sjzar/ips/format"
	"github.com/sjzar/ips/internal/operate"
	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

// StandardReader is a Reader designed for IP databases. Unlike the Reader in the format package,
// this reader processes the IP information it reads, such as filtering fields, rewriting data,
// supplementing data, etc., controlled by operate.IPOperateChain.
type StandardReader struct {
	DBReader     format.Reader
	OperateChain *operate.IPOperateChain

	meta *model.Meta
}

// NewStandardReader initializes and returns a new StandardReader.
func NewStandardReader(dbReader format.Reader, operateChain *operate.IPOperateChain) *StandardReader {
	if operateChain == nil {
		operateChain = operate.NewIPOperateChain()
	}
	return &StandardReader{
		DBReader:     dbReader,
		OperateChain: operateChain,
	}
}

// Meta returns the meta-information of the IP database.
func (s *StandardReader) Meta() *model.Meta {
	if s.meta != nil {
		return s.meta
	}

	meta := s.DBReader.Meta()
	s.meta = &model.Meta{
		MetaVersion: meta.MetaVersion,
		Format:      meta.Format,
		IPVersion:   meta.IPVersion,
		Fields:      meta.Fields,
		FieldAlias:  meta.FieldAlias,
	}
	return s.meta
}

// Find retrieves IP information based on the given IP address.
func (s *StandardReader) Find(ip net.IP) (*model.IPInfo, error) {
	info, err := s.DBReader.Find(ip)
	if err != nil {
		return nil, err
	}

	if s.OperateChain != nil {
		if err := s.OperateChain.Do(info); err != nil {
			return nil, err
		}
	}

	return info, nil
}

type StandardReaderOption struct {
	IPVersion int

	Fields []string
}

// SetOption configures the StandardReader with the provided option.
func (s *StandardReader) SetOption(option interface{}) error {
	if opt, ok := option.(StandardReaderOption); ok {
		if opt.IPVersion > 0 {
			if err := s.setIPVersion(opt.IPVersion); err != nil {
				return err
			}
		}
		if len(opt.Fields) > 0 {
			if err := s.setFields(opt.Fields); err != nil {
				return err
			}
		}
	}

	return nil
}

// Close closes the IP database.
func (s *StandardReader) Close() error {
	return s.DBReader.Close()
}

// setIPVersion filters the IP data based on the IP version.
// For example, if the database contains both IPv4 and IPv6 data but only IPv4 data is needed.
func (s *StandardReader) setIPVersion(ipVersion int) error {
	if ipVersion == 0 {
		return nil
	}

	meta := s.Meta()
	if ipVersion&meta.IPVersion != ipVersion {
		return errors.ErrUnsupportedIPVersion
	}

	meta.IPVersion = ipVersion
	s.meta = meta
	return nil
}

// setFields filters the fields of the IP data.
func (s *StandardReader) setFields(fields []string) error {
	if len(fields) == 0 {
		return nil
	}

	meta := s.Meta()
	meta.Fields = fields
	s.meta = meta
	return nil
}
