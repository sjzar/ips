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

package plain

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"strings"

	"github.com/sjzar/ips/ipnet"
	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

const (
	DBFormat   = "plain"
	DBExt      = ".txt"
	MetaPrefix = "# Meta: "
	FieldSep   = ","
)

// Reader is a structure that provides functionalities to read from Plain Text.
type Reader struct {
	file    string
	meta    *model.Meta
	fd      *os.File
	scanner *bufio.Scanner
	done    bool
}

// NewReader initializes a new instance of Reader.
func NewReader(file string) (*Reader, error) {

	meta, err := ParseMeta(file)
	if err != nil {
		return nil, err
	}

	meta.Format = DBFormat

	return &Reader{
		file: file,
		meta: meta,
	}, nil
}

// Find retrieves IP information based on the given IP address.
// FIXME this method returns the next available IP info sequentially from the file.
func (r *Reader) Find(ip net.IP) (*model.IPInfo, error) {
	return r.Next()
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
	if r.fd != nil {
		return r.fd.Close()
	}
	return nil
}

// ParseMeta reads and parses the metadata from the beginning of the plain text IP database file.
func ParseMeta(file string) (*model.Meta, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	ret := &model.Meta{}
	success := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, MetaPrefix) {
			success = true
			metaStr := line[len(MetaPrefix):]
			if err := json.Unmarshal([]byte(metaStr), ret); err != nil {
				return nil, err
			}
			break
		}
	}

	if !success {
		return nil, errors.ErrMetaMissing
	}

	return ret, nil
}

// Next retrieves the next IP information from the file.
func (r *Reader) Next() (*model.IPInfo, error) {
	if r.scanner == nil {
		var err error
		r.fd, err = os.Open(r.file)
		if err != nil {
			return nil, err
		}
		r.scanner = bufio.NewScanner(r.fd)
	}

	if r.done {
		return nil, errors.ErrReadCompleted
	}

	ret := &model.IPInfo{
		Fields: r.meta.Fields,
		Data:   make(map[string]string),
	}
	for r.scanner.Scan() {
		line := r.scanner.Text()
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		split := strings.SplitN(line, "\t", 2)
		if len(split) != 2 {
			continue
		}

		_, ipNet, err := net.ParseCIDR(split[0])
		if err != nil {
			continue
		}

		values := strings.Split(split[1], FieldSep)
		if len(values) != len(r.meta.Fields) {
			continue
		}

		ret.IP = ipNet.IP
		ret.IPNet = ipnet.NewRange(ipNet)
		for i, field := range r.meta.Fields {
			ret.Data[field] = values[i]
		}

		if ipnet.IsLastIP(ret.IPNet.End, r.meta.IsIPv6Support()) {
			r.done = true
		}

		break
	}

	if err := r.scanner.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}
