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

package ip2region

import "github.com/sjzar/ips/model"

// 中国|0|福建省|福州市|电信
// 澳大利亚|0|维多利亚|墨尔本|0
// 泰国|0|清莱府|0|TOT

const (
	// FieldCountry 国家
	FieldCountry = "country"

	// FieldRegion 区域
	// 例如 华东/华北 等，在v2中基本为空
	FieldRegion = "region"

	// FieldProvince 省份
	FieldProvince = "province"

	// FieldCity 城市
	FieldCity = "city"

	// FieldISP 运营商
	FieldISP = "isp"
)

// FullFields 全字段列表
var FullFields = []string{
	FieldCountry,
	FieldRegion,
	FieldProvince,
	FieldCity,
	FieldISP,
}

// CommonFieldsMap 公共字段映射
var CommonFieldsMap = map[string]string{
	model.Country:  FieldCountry,
	model.Province: FieldProvince,
	model.City:     FieldCity,
	model.ISP:      FieldISP,
}

// FieldsFormat 字段格式化，并补充公共字段
func FieldsFormat(data map[string]string) map[string]string {

	// Fill Common Fields
	for k, v := range CommonFieldsMap {
		data[k] = data[v]
	}

	return data
}
