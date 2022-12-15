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
	ipv4once.Do(getIPv4)
	return ipv4
}

func getIPv4() {
	format := rootDBFormat
	if len(format) == 0 {
		format = conf.Conf.IPv4Format
	}
	file := rootDBFile
	if len(file) == 0 {
		file = conf.GetIPv4File()
	}
	database, err := db.NewDatabase(format, file)
	if err != nil {
		if file == conf.GetIPv4File() {
			update()
			database, err = db.NewDatabase(format, file)
		}
		if err != nil {
			log.Println("read database failed", file, err)
			os.Exit(1)
		}
	}
	var fields []string
	if len(conf.Conf.Fields) != 0 {
		fields = strings.Split(conf.Conf.Fields, ",")
	}
	if len(rootFields) != 0 {
		fields = strings.Split(rootFields, ",")
	}
	if len(fields) == 0 {
		fields = database.Meta().Fields
	}
	selector := ipio.NewFieldSelector(strings.Join(fields, ","))
	ipv4 = ipio.NewDBScanner(database, selector, nil)
}

// GetIPv6 returns a ipio.Reader for IPv6
func GetIPv6() ipio.Reader {
	ipv6once.Do(getIPv6)
	return ipv6
}

func getIPv6() {
	format := rootDBFormat
	if len(format) == 0 {
		format = conf.Conf.IPv4Format
	}
	file := rootDBFile
	if len(file) == 0 {
		file = conf.GetIPv6File()
	}
	database, err := db.NewDatabase(format, file)
	if err != nil {
		if file == conf.GetIPv6File() {
			update()
			database, err = db.NewDatabase(format, file)
		}
		if err != nil {
			log.Println("read database failed", file, err)
			os.Exit(1)
		}
	}
	var fields []string
	if len(conf.Conf.Fields) != 0 {
		fields = strings.Split(conf.Conf.Fields, ",")
	}
	if len(rootFields) != 0 {
		fields = strings.Split(rootFields, ",")
	}
	if len(fields) == 0 {
		fields = database.Meta().Fields
	}
	selector := ipio.NewFieldSelector(strings.Join(fields, ","))
	ipv6 = ipio.NewDBScanner(database, selector, nil)
}
