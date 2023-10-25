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
	"bytes"
	"encoding/json"
	"net"
	"os"
	"strings"

	"github.com/sjzar/ips/format/ipdb"
	"github.com/sjzar/ips/format/ipdb/sdk"
	"github.com/sjzar/ips/ipnet"
	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

const (
	DBFormat   = "plain"
	DBExt      = ".txt"
	MetaPrefix = "# Meta: "
	FieldSep   = ","
	FieldData  = "text"
)

// Reader is a structure that provides functionalities to read from Plain Text.
type Reader struct {
	file string
	meta *model.Meta
	db   *sdk.City
}

// NewReader initializes a new instance of Reader.
func NewReader(file string) (*Reader, error) {

	meta, db, err := Load(file)
	if err != nil {
		return nil, err
	}

	meta.Format = DBFormat

	return &Reader{
		file: file,
		meta: meta,
		db:   db,
	}, nil
}

// Find retrieves IP information based on the given IP address.
func (r *Reader) Find(ip net.IP) (*model.IPInfo, error) {
	data, ipNet, err := r.db.FindMap(ip.String(), "CN")
	if err != nil {
		return nil, err
	}
	values := strings.Split(data[FieldData], FieldSep)

	ret := &model.IPInfo{
		IP:     ip,
		IPNet:  ipnet.NewRange(ipNet),
		Fields: r.meta.Fields,
		Data:   make(map[string]string),
	}
	for i, field := range r.meta.Fields {
		ret.Data[field] = values[i]
	}
	ret.AddCommonFieldAlias(r.meta.FieldAlias)

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

// Load opens the specified file, reads its contents, and initializes an IP database.
// The function first extracts meta information from the file, then parses the IP data,
// and finally creates an IP database based on the parsed data.
// FIXME: The plain database format reader has not been implemented yet, so I made a wrap in ipdb format.
func Load(file string) (*model.Meta, *sdk.City, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	// Extract meta information from the file.
	scanner := bufio.NewScanner(f)
	meta, err := extractMetaInfo(scanner)
	if err != nil {
		return nil, nil, err
	}

	// Create a new IP database writer with the extracted meta information.
	writer, err := ipdb.NewWriter(&model.Meta{
		IPVersion:  meta.IPVersion,
		Fields:     []string{FieldData},
		FieldAlias: map[string]string{},
	})
	if err != nil {
		return nil, nil, err
	}

	// Parse the IP data from the file and load it into the writer.
	err = parseIPData(scanner, writer, meta.IsIPv6Support())
	if err != nil {
		return nil, nil, err
	}

	// Initialize the IP database using the loaded data.
	db, err := initializeDatabase(writer)
	if err != nil {
		return nil, nil, err
	}

	return meta, db, nil
}

// extractMetaInfo reads the meta information from the file and returns it.
func extractMetaInfo(scanner *bufio.Scanner) (*model.Meta, error) {
	meta := &model.Meta{}
	success := false

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, MetaPrefix) {
			success = true
			metaStr := line[len(MetaPrefix):]
			if err := json.Unmarshal([]byte(metaStr), meta); err != nil {
				return nil, err
			}
			break
		}
	}

	if !success {
		return nil, errors.ErrMetaMissing
	}

	return meta, nil
}

// parseIPData processes the IP information from the file and loads it into the provided writer.
func parseIPData(scanner *bufio.Scanner, writer *ipdb.Writer, isIPv6Support bool) error {
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
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

		info := &model.IPInfo{
			IP:     ipNet.IP,
			IPNet:  ipnet.NewRange(ipNet),
			Data:   map[string]string{FieldData: split[1]},
			Fields: []string{FieldData},
		}

		if err := writer.Insert(info); err != nil {
			return err
		}

		if ipnet.IsLastIP(info.IPNet.End, isIPv6Support) {
			break
		}
	}
	return nil
}

// initializeDatabase creates the IP database and returns it.
func initializeDatabase(writer *ipdb.Writer) (*sdk.City, error) {
	r := &bytes.Buffer{}
	if _, err := writer.WriteTo(r); err != nil {
		return nil, err
	}

	db, err := sdk.NewCityByIO(r)
	if err != nil {
		return nil, err
	}

	return db, nil
}
