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

package awdb

import (
	"github.com/sjzar/ips/model"
)

// https://mall.ipplus360.com/pros/IPVFourGeoDB
// Example:
// country:中国
// province:浙江省
// city:绍兴市
// isp:中国电信
// continent:亚洲
// timezone:UTC+8
// latwgs:29.998742
// lngwgs:120.581963
// adcode:330600
// accuracy:城市
// areacode:CN
// asnumber:4134
// owner:中国电信
// radius:71.2163
// source:数据挖掘
// zipcode:131200

const (
	// FieldCountry 国家
	FieldCountry = "country"

	// FieldProvince 省份
	FieldProvince = "province"

	// FieldCity 城市
	FieldCity = "city"

	// FieldISP 运营商
	FieldISP = "isp"

	// FieldContinent 大洲
	FieldContinent = "continent"

	// FieldTimeZone 时区
	FieldTimeZone = "timezone"

	// FieldLatwgs WGS84坐标系纬度
	FieldLatwgs = "latwgs"

	// FieldLngwgs WGS84坐标系经度
	FieldLngwgs = "lngwgs"

	// FieldAdcode 行政区划代码
	FieldAdcode = "adcode"

	// FieldAccuracy 定位精度
	FieldAccuracy = "accuracy"

	// FieldAreaCode 国家编码
	FieldAreaCode = "areacode"

	// FieldASNumber 自治域编码
	FieldASNumber = "asnumber"

	// FieldOwner 所属机构
	FieldOwner = "owner"

	// FieldRadius 定位半径
	FieldRadius = "radius"

	// FieldSource 定位方式
	FieldSource = "source"

	// FieldZipCode 邮编
	FieldZipCode = "zipcode"
)

// FullFields 全字段列表
var FullFields = []string{
	FieldCountry,
	FieldProvince,
	FieldCity,
	FieldISP,
	FieldContinent,
	FieldTimeZone,
	FieldLatwgs,
	FieldLngwgs,
	FieldAdcode,
	FieldAccuracy,
	FieldAreaCode,
	FieldASNumber,
	FieldOwner,
	FieldRadius,
	FieldSource,
	FieldZipCode,
}

// CommonFieldsMap 公共字段映射
var CommonFieldsMap = map[string]string{
	model.Country:        FieldCountry,
	model.Province:       FieldProvince,
	model.City:           FieldCity,
	model.ISP:            FieldISP,
	model.Continent:      FieldContinent,
	model.UTCOffset:      FieldTimeZone,
	model.Latitude:       FieldLatwgs,
	model.Longitude:      FieldLngwgs,
	model.ChinaAdminCode: FieldAdcode,
}

// FieldsFormat 字段格式化，并补充公共字段
func FieldsFormat(record map[string]interface{}) map[string]string {
	data := make(map[string]string)

	for k, v := range record {
		data[k] = string(v.([]byte))
	}

	// Fill Common Fields
	for k, v := range CommonFieldsMap {
		data[k] = data[v]
	}

	return data
}
