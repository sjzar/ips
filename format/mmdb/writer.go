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
	"io"
	"net"
	"strings"

	"github.com/maxmind/mmdbwriter"

	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

// Writer provides functionalities to write IP data into MMDB format.
type Writer struct {
	meta   *model.Meta      // Metadata for the IP database
	writer *mmdbwriter.Tree // MMDB writer instance
}

// NewWriter initializes a new Writer instance for writing IP data in MMDB format.
func NewWriter(meta *model.Meta) (*Writer, error) {

	opts := mmdbwriter.Options{
		DatabaseType: "GeoIP2-City",
	}

	writer, err := mmdbwriter.New(opts)
	if err != nil {
		return nil, err
	}

	return &Writer{
		meta:   meta,
		writer: writer,
	}, nil
}

// WriterOption provides options for the Writer.
type WriterOption struct {
}

// SetOption sets the provided options to the Writer.
// Currently, it supports mmdbwriter.Options for the MMDB writer.
func (w *Writer) SetOption(option interface{}) error {
	if _, ok := option.(WriterOption); ok {
		return nil
	}
	if opts, ok := option.(mmdbwriter.Options); ok {
		writer, err := mmdbwriter.New(opts)
		if err != nil {
			return err
		}
		w.writer = writer
	}
	return nil
}

// Insert adds the given IP information into the writer.
func (w *Writer) Insert(info *model.IPInfo) error {

	fields := model.ConvertToDBFields(w.meta.Fields, w.meta.FieldAlias, CommonFieldsAlias)
	values := info.Values()
	if len(values) != len(fields) {
		return errors.ErrMismatchedFieldsLength
	}
	data := ConvertMap(fields, values)

	for _, ipNet := range info.IPNet.IPNets() {
		_, network, err := net.ParseCIDR(ipNet.String())
		if err != nil || network == nil {
			continue
		}
		if err := w.writer.Insert(network, data); err != nil {
			if strings.Contains(err.Error(), "which is in a reserved network") ||
				strings.Contains(err.Error(), "which is in an aliased network") {
				err = nil
			}
			return err
		}
	}
	return nil
}

// WriteTo writes the IP data into the provided writer in MMDB format.
func (w *Writer) WriteTo(iw io.Writer) (int64, error) {
	return w.writer.WriteTo(iw)
}
