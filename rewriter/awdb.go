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

package rewriter

import (
	"log"
	"net"

	"github.com/sjzar/ips/db/awdb"
	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

const (
	AWDBFieldLine = "line"
)

// AWDBIPLineRewriter AWDB数据库的IP线路改写
// 根据线路数据库匹配 ISP
type AWDBIPLineRewriter struct {
	awdb *awdb.Database
	rw   Rewriter
}

// NewIPLineRewriter 初始化实例
func NewIPLineRewriter(ipLineFile string, rw Rewriter) (Rewriter, error) {
	if rw == nil {
		rw = DefaultRewriter
	}

	db, err := awdb.New(ipLineFile)
	if err != nil {
		log.Println("new awdb database failed", err)
		return nil, err
	}
	return &AWDBIPLineRewriter{
		awdb: db,
		rw:   rw,
	}, nil
}

// Rewrite 改写
func (r *AWDBIPLineRewriter) Rewrite(ip net.IP, ipRange *ipx.Range, data map[string]string) (net.IP, *ipx.Range, map[string]string, error) {
	ipRangeLine, dataLine, err := r.awdb.Find(ip)
	if err != nil {
		return nil, nil, nil, err
	}

	if line, ok := dataLine[AWDBFieldLine]; ok {
		if ok := ipRange.CommonRange(ip, ipRangeLine); ok {
			data[model.ISP] = line
		}
	}

	return r.rw.Rewrite(ip, ipRange, data)
}
