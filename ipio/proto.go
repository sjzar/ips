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

package ipio

import (
	"io"
	"net"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

// Reader IP库读取工具
// 从数据库中获取IP的地理位置信息，并格式化数据
type Reader interface {

	// Meta 返回元数据
	Meta() model.Meta

	// Find 查询IP所在网段和对应数据
	Find(ip net.IP) (*ipx.Range, []string, error)

	// Close 关闭
	Close() error
}

// Scanner IP库扫描工具
type Scanner interface {

	// Meta 返回元数据
	Meta() model.Meta

	// Scan 扫描下一条数据，返回 true 表示有数据，false 表示没有数据
	Scan() bool

	// Result 返回结果
	Result() (*ipx.Range, []string)

	// Err 返回错误
	Err() error

	// Close 关闭
	Close() error
}

// Writer IP库写入工具
type Writer interface {

	// Insert 插入数据
	Insert(ipr *ipx.Range, values []string) error

	// Save 保存数据
	Save(w io.Writer) error
}
