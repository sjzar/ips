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
	"github.com/sjzar/ips/db/awdb"
	"github.com/sjzar/ips/db/ip2region"
	"github.com/sjzar/ips/db/ipdb"
	"github.com/sjzar/ips/db/mmdb"
	"github.com/sjzar/ips/db/qqwry"
	"github.com/sjzar/ips/db/zxinc"
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

func init() {
	register(FormatIPDB, ".ipdb", []string{}, func(s string) (Database, error) { return ipdb.New(s) })
	register(FormatAWDB, ".awdb", []string{}, func(s string) (Database, error) { return awdb.New(s) })
	register(FormatMMDB, ".mmdb", []string{"GeoLite2-City", "dbip-city-lite"}, func(s string) (Database, error) { return mmdb.New(s) })
	register(FormatQQWry, ".dat", []string{"qqwry"}, func(s string) (Database, error) { return qqwry.New(s) })
	register(FormatZXInc, ".db", []string{"zxipv6wry"}, func(s string) (Database, error) { return zxinc.New(s) })
	register(FormatIP2Region, ".xdb", []string{"ip2region"}, func(s string) (Database, error) { return ip2region.New(s) })
}

var (
	// Formats 支持的数据格式
	Formats = map[string]func(string) (Database, error){}

	// Exts 支持的文件扩展名
	Exts = map[string]func(string) (Database, error){}

	// CommonNames 支持的常用文件名称
	CommonNames = map[string]func(string) (Database, error){}
)

// register 注册可支持的数据格式
func register(format, ext string, commonName []string, fn func(string) (Database, error)) {

	if len(format) > 0 {
		if _, ok := Formats[format]; ok {
			panic("ips: format already registered " + format)
		}
		Formats[format] = fn
	}

	if len(ext) > 0 {
		if _, ok := Exts[ext]; ok {
			panic("ips: ext already registered " + ext)
		}
		Exts[ext] = fn
	}

	if len(commonName) > 0 {
		for _, name := range commonName {
			if len(name) == 0 {
				continue
			}
			if _, ok := CommonNames[name]; ok {
				panic("ips: common name already registered " + name)
			}
			CommonNames[name] = fn
		}
	}
}
