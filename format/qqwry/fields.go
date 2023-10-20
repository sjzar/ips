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

package qqwry

import (
	"github.com/sjzar/ips/pkg/model"
)

// cidr,country,area
// 1.2.2.0/24,北京市海淀区,北龙中网(北京)科技有限公司
// 1.2.8.0/24,中国,网络信息中心
// 1.2.9.0/24,广东省,电信
// 1.8.151.0/24,香港,特别行政区
// 1.10.140.0/25,泰国, CZ88.NET
// 1.12.36.0/22,广东省广州市,腾讯云
// * Country 和 Area 中的数据比较混乱，需要后续改写处理

const (
	FieldCountry = "country"
	FieldArea    = "area"
)

// FullFields 全字段列表
var FullFields = []string{
	FieldCountry,
	FieldArea,
}

// CommonFieldsAlias 公共字段到数据库字段映射
var CommonFieldsAlias = map[string]string{
	model.Country: FieldCountry,
	model.ISP:     FieldArea,
}
