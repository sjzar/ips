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

package zxinc

import (
	"github.com/sjzar/ips/pkg/model"
)

// cidr,country,area
// 2001:200:120::/48,日本    东京都  品川区,Sony Computer Science 研究所
// 2001:200:178::/48,美国    California州    San José,Sony NCSA实验室
// 2001:250:1::/48,中国    北京市,教育网(CERNET)网络运行部
// ffff:ffff:ffff:fff0::/60,,ZX公网IPv6库 20210511版
// * Country中有多级数据，需要处理

const (
	FieldCountry = "country"
	FieldArea    = "area"
)

// FullFields 全字段列表
var FullFields = []string{
	FieldCountry,
	FieldArea,
}

// CommonFieldsAlias 公共字段映射
var CommonFieldsAlias = map[string]string{
	model.Country: FieldCountry,
	model.ISP:     FieldArea,
}
