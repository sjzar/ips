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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

// Writer provides functionalities to write IP data into plain text format.
type Writer struct {
	meta   *model.Meta
	iw     io.Writer
	buffer *bytes.Buffer
}

// NewWriter initializes a new Writer instance for writing IP data in plain text format.
// If the io.Writer is not nil, it will directly write data into it.
func NewWriter(meta *model.Meta) (*Writer, error) {
	ret := &Writer{
		meta: meta,
	}

	return ret, nil
}

// WriterOption provides options for the Writer.
type WriterOption struct {
	// IW is for immediate output to the provided writer.
	IW io.Writer
}

// SetOption sets the provided options to the Writer.
func (w *Writer) SetOption(option interface{}) error {
	if opt, ok := option.(WriterOption); ok {
		if opt.IW != nil {
			w.iw = opt.IW
			if err := w.Header(); err != nil {
				return err
			}
		}
		return nil
	}

	return nil
}

// Insert adds the given IP information into the writer.
func (w *Writer) Insert(info *model.IPInfo) error {
	if w.iw == nil {
		w.buffer = bytes.NewBuffer([]byte{})
		w.iw = w.buffer

		if err := w.Header(); err != nil {
			return err
		}
	}

	for _, ipNet := range info.IPNet.IPNets() {
		if _, err := fmt.Fprintf(w.iw, "%s\t%s\n", ipNet.String(), strings.Join(info.Values(), FieldSep)); err != nil {
			return err
		}
	}

	return nil
}

// WriteTo writes the buffered data into the provided writer.
func (w *Writer) WriteTo(writer io.Writer) (int64, error) {
	if w.buffer == nil {
		return 0, nil
	}

	return w.buffer.WriteTo(writer)
}

// Header writes the header information for the IP database.
func (w *Writer) Header() error {
	if w.iw == nil {
		return errors.ErrNilWriter
	}

	str := fmt.Sprintf("# Dump Time: %s\n", time.Now().Local().Format("2006-01-02 15:04:05"))
	str += fmt.Sprintf("# Fields: %s\n", strings.Join(w.meta.Fields, FieldSep))
	str += fmt.Sprintf("# IP Version: %d\n", w.meta.IPVersion)
	b, _ := json.Marshal(w.meta)
	str += fmt.Sprintf("%s%s\n", MetaPrefix, string(b))
	if _, err := fmt.Fprint(w.iw, str); err != nil {
		return err
	}

	return nil
}
