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

package format

import (
	"io"
	"path/filepath"

	"github.com/sjzar/ips/format/ipdb"
	"github.com/sjzar/ips/format/mmdb"
	"github.com/sjzar/ips/format/plain"
	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

// Writer defines an interface for writing IP databases.
type Writer interface {

	// SetOption configures the Writer with the provided option.
	SetOption(option interface{}) error

	// Insert IP information
	Insert(info *model.IPInfo) error

	// WriteTo writes data to io.Writer
	WriteTo(w io.Writer) (int64, error)

	// WriterFormat returns the format of the writer.
	WriterFormat() string
}

// NewWriter creates a Writer based on its format or file name.
func NewWriter(format string, file string, meta *model.Meta) (Writer, error) {
	if fn, ok := WriterFormats[format]; ok {
		return fn(meta)
	}

	if fn, ok := WriterExts[filepath.Ext(file)]; ok {
		return fn(meta)
	}

	return nil, errors.ErrUnsupportedFormat
}

var (
	WriterFormats = map[string]func(meta *model.Meta) (Writer, error){
		ipdb.DBFormat:  func(meta *model.Meta) (Writer, error) { return ipdb.NewWriter(meta) },
		mmdb.DBFormat:  func(meta *model.Meta) (Writer, error) { return mmdb.NewWriter(meta) },
		plain.DBFormat: func(meta *model.Meta) (Writer, error) { return plain.NewWriter(meta) },
	}
	WriterExts = map[string]func(meta *model.Meta) (Writer, error){
		ipdb.DBExt:  func(meta *model.Meta) (Writer, error) { return ipdb.NewWriter(meta) },
		mmdb.DBExt:  func(meta *model.Meta) (Writer, error) { return mmdb.NewWriter(meta) },
		plain.DBExt: func(meta *model.Meta) (Writer, error) { return plain.NewWriter(meta) },
	}
)

// registerWriter is a helper function to register a writer to the provided map.
func registerWriter(m map[string]func(meta *model.Meta) (Writer, error), key string, fn func(meta *model.Meta) (Writer, error)) {
	mu.Lock()
	defer mu.Unlock()
	m[key] = fn
}

// RegisterWriterFormat registers a Writer by its format.
func RegisterWriterFormat(name string, fn func(meta *model.Meta) (Writer, error)) {
	if name == "" || fn == nil {
		return
	}
	registerWriter(WriterFormats, name, fn)
}

// RegisterWriterExt registers a Writer by its file extension.
func RegisterWriterExt(ext string, fn func(meta *model.Meta) (Writer, error)) {
	if ext == "" || fn == nil {
		return
	}
	registerWriter(WriterExts, ext, fn)
}
