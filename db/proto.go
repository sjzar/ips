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

	"github.com/sjzar/ips/db/awdb"
	"github.com/sjzar/ips/db/ip2region"
	"github.com/sjzar/ips/db/ipdb"
	"github.com/sjzar/ips/db/mmdb"
	"github.com/sjzar/ips/db/qqwry"
	"github.com/sjzar/ips/db/zxinc"
	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

const (
	// FormatIPDB IPDB格式
	// Official: https://www.ipip.net/
	FormatIPDB = "ipdb"

	// FormatAWDB AWDB格式
	// Official: https://www.ipplus360.com/
	FormatAWDB = "awdb"

	// FormatMMDB MMDB格式
	// Official: https://www.maxmind.com/
	FormatMMDB = "mmdb"

	// FormatQQWry QQWry格式
	// Official: https://www.cz88.net/
	FormatQQWry = "qqwry"

	// FormatZXInc ZXInc格式
	// Official: https://ip.zxinc.org/
	FormatZXInc = "zxinc"

	// FormatIP2Region (v2)
	// Official: https://github.com/lionsoul2014/ip2region
	FormatIP2Region = "ip2region"
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
		if _file, err := filepath.Abs(file); err != nil {
			file = _file
		}
	}

	if len(format) == 0 {
		switch {
		case filepath.Ext(file) == ".ipdb":
			format = FormatIPDB
		case filepath.Ext(file) == ".awdb":
			format = FormatAWDB
		case filepath.Ext(file) == ".mmdb":
			format = FormatMMDB
		case filepath.Ext(file) == ".xdb":
			format = FormatIP2Region
		case strings.HasSuffix(file, "qqwry.dat"):
			format = FormatQQWry
		case strings.HasSuffix(file, "zxipv6wry.db"):
			format = FormatZXInc
		}
	}

	switch format {
	case FormatIPDB:
		return ipdb.New(file)
	case FormatAWDB:
		return awdb.New(file)
	case FormatQQWry:
		return qqwry.New(file)
	case FormatMMDB:
		return mmdb.New(file)
	case FormatZXInc:
		return zxinc.New(file)
	case FormatIP2Region:
		return ip2region.New(file)
	}

	return nil, errors.ErrDBFormatNotSupported
}
