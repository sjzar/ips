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

package ips

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sjzar/ips/db"
	"github.com/sjzar/ips/ipio"

	"github.com/sjzar/ips/cmd/ips/conf"
)

var (
	// IPv4
	ipv4     ipio.Reader
	ipv4once sync.Once

	// IPv6
	ipv6     ipio.Reader
	ipv6once sync.Once
)

// GetIPv4 returns a ipio.Reader for IPv4
func GetIPv4() ipio.Reader {
	ipv4once.Do(func() {
		format := rootDBFormat
		if len(format) == 0 {
			format = conf.Conf.IPv4Format
		}
		file := rootDBFile
		if len(file) == 0 {
			file = conf.ConfigPath + "/" + conf.Conf.IPv4File
		}
		database, err := db.NewDatabase(format, file)
		if err != nil {
			log.Println("read database failed", file, err)
			if file == conf.ConfigPath+"/"+conf.Conf.IPv4File {
				log.Println("use ips update first")
			}
			os.Exit(1)
		}
		fields := conf.Conf.Fields
		if len(rootFields) != 0 {
			fields = strings.Split(rootFields, ",")
		}
		if len(fields) == 0 {
			fields = database.Meta().Fields
		}
		selector := ipio.NewFieldSelector(strings.Join(fields, ","))
		ipv4 = ipio.NewDBScanner(database, selector, nil)
	})

	return ipv4
}

// GetIPv6 returns a ipio.Reader for IPv6
func GetIPv6() ipio.Reader {
	ipv6once.Do(func() {
		format := rootDBFormat
		if len(format) == 0 {
			format = conf.Conf.IPv4Format
		}
		file := rootDBFile
		if len(file) == 0 {
			file = conf.ConfigPath + "/" + conf.Conf.IPv6File
		}
		database, err := db.NewDatabase(format, file)
		if err != nil {
			log.Println("read database failed", file, err)
			if file == conf.ConfigPath+"/"+conf.Conf.IPv6File {
				log.Println("use ips update first")
			}
			os.Exit(1)
		}
		fields := conf.Conf.Fields
		if len(rootFields) != 0 {
			fields = strings.Split(rootFields, ",")
		}
		if len(fields) == 0 {
			fields = database.Meta().Fields
		}
		selector := ipio.NewFieldSelector(strings.Join(fields, ","))
		ipv6 = ipio.NewDBScanner(database, selector, nil)
	})

	return ipv6
}
