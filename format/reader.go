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
	"net"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sjzar/ips/format/awdb"
	"github.com/sjzar/ips/format/ip2region"
	"github.com/sjzar/ips/format/ipdb"
	"github.com/sjzar/ips/format/mmdb"
	"github.com/sjzar/ips/format/plain"
	"github.com/sjzar/ips/format/qqwry"
	"github.com/sjzar/ips/format/zxinc"
	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

// Reader defines an interface for reading IP databases.
type Reader interface {

	// Meta returns the meta-information of the IP database.
	Meta() *model.Meta

	// Find retrieves IP information based on the given IP address.
	Find(ip net.IP) (*model.IPInfo, error)

	// SetOption configures the Reader with the provided option.
	SetOption(option interface{}) error

	// Close closes the IP database.
	Close() error
}

// NewReader creates a Reader based on its format or file name.
func NewReader(format, file string) (Reader, error) {
	if fn, ok := ReaderFormats[format]; ok {
		return fn(file)
	}

	if fn, ok := ReaderExts[filepath.Ext(file)]; ok {
		return fn(file)
	}

	for commonName, fn := range ReaderCommonNames {
		if strings.HasPrefix(filepath.Base(file), commonName) {
			return fn(file)
		}
	}

	return nil, errors.ErrUnsupportedFormat
}

var (
	mu            sync.Mutex
	ReaderFormats = map[string]func(string) (Reader, error){
		awdb.DBFormat:      func(file string) (Reader, error) { return awdb.NewReader(file) },
		ip2region.DBFormat: func(file string) (Reader, error) { return ip2region.NewReader(file) },
		ipdb.DBFormat:      func(file string) (Reader, error) { return ipdb.NewReader(file) },
		mmdb.DBFormat:      func(file string) (Reader, error) { return mmdb.NewReader(file) },
		plain.DBFormat:     func(file string) (Reader, error) { return plain.NewReader(file) },
		qqwry.DBFormat:     func(file string) (Reader, error) { return qqwry.NewReader(file) },
		zxinc.DBFormat:     func(file string) (Reader, error) { return zxinc.NewReader(file) },
	}
	ReaderExts = map[string]func(string) (Reader, error){
		awdb.DBExt:      func(file string) (Reader, error) { return awdb.NewReader(file) },
		ip2region.DBExt: func(file string) (Reader, error) { return ip2region.NewReader(file) },
		ipdb.DBExt:      func(file string) (Reader, error) { return ipdb.NewReader(file) },
		mmdb.DBExt:      func(file string) (Reader, error) { return mmdb.NewReader(file) },
		plain.DBExt:     func(file string) (Reader, error) { return plain.NewReader(file) },
		qqwry.DBExt:     func(file string) (Reader, error) { return qqwry.NewReader(file) },
		zxinc.DBExt:     func(file string) (Reader, error) { return zxinc.NewReader(file) },
	}
	ReaderCommonNames = map[string]func(string) (Reader, error){}
)

// registerReader is a helper function to register a reader to the provided map.
func registerReader(m map[string]func(string) (Reader, error), key string, fn func(string) (Reader, error)) {
	mu.Lock()
	defer mu.Unlock()
	m[key] = fn
}

// RegisterReaderFormat registers a Reader by its format.
func RegisterReaderFormat(name string, fn func(string) (Reader, error)) {
	if name == "" || fn == nil {
		return
	}
	registerReader(ReaderFormats, name, fn)
}

// RegisterReaderExt registers a Reader by its file extension.
func RegisterReaderExt(ext string, fn func(string) (Reader, error)) {
	if ext == "" || fn == nil {
		return
	}
	registerReader(ReaderExts, ext, fn)
}

// RegisterReaderCommonName registers a Reader by its common file name.
func RegisterReaderCommonName(name string, fn func(string) (Reader, error)) {
	if name == "" || fn == nil {
		return
	}
	registerReader(ReaderCommonNames, name, fn)
}
