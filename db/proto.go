/*
 * Copyright (c) 2022 shenjunzheng@gmail.com
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

package db

import (
	"net"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

const (
	// FormatSep 格式分隔符
	FormatSep = ":"
)

// Database IP 数据库
type Database interface {

	// Meta 返回元数据
	Meta() model.Meta

	// Find 查询 IP 对应的网段和结果
	Find(ip net.IP) (*ipx.Range, map[string]string, error)

	// Close 关闭数据库实例
	Close() error
}

// NewDatabase 初始化 IP 数据库
func NewDatabase(format, file string) (Database, error) {

	// support format:file
	// ipdb:./data/ipip.ipdb
	// awdb:./data/ipplus.awdb
	if len(format) == 0 && strings.Contains(file, FormatSep) {
		split := strings.SplitN(file, FormatSep, 2)
		format = split[0]
		file = split[1]
		if strings.Contains(file, "~") {
			if u, err := user.Current(); err == nil {
				file = strings.Replace(file, "~", u.HomeDir, -1)
			}
		}
		if _file, err := filepath.Abs(file); err == nil {
			file = _file
		}
	}

	if fn, ok := Formats[format]; ok {
		return fn(file)
	}

	if fn, ok := Exts[filepath.Ext(file)]; ok {
		return fn(file)
	}

	for commonName, fn := range CommonNames {
		if strings.HasPrefix(filepath.Base(file), commonName) {
			return fn(file)
		}
	}

	return nil, errors.ErrDBFormatNotSupported
}
