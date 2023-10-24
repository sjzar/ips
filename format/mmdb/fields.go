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

package mmdb

import (
	"github.com/sjzar/ips/pkg/model"
)

const (

	// FieldCity 城市
	FieldCity = "city"

	// FieldContinent 大洲
	FieldContinent = "continent"

	// FieldCountry 国家
	FieldCountry = "country"

	// FieldSubdivisions 行政区
	FieldSubdivisions = "subdivisions"

	// FieldAccuracyRadius 定位精度
	FieldAccuracyRadius = "accuracy_radius"

	// FieldLatitude 纬度
	FieldLatitude = "latitude"

	// FieldLongitude 经度
	FieldLongitude = "longitude"

	// FieldMetroCode 城市代码
	FieldMetroCode = "metro_code"

	// FieldTimeZone 时区
	FieldTimeZone = "time_zone"

	// FieldPostalCode 邮政编码
	FieldPostalCode = "postal_code"

	// FieldRegisteredCountry 注册国家
	FieldRegisteredCountry = "registered_country"

	// FieldRepresentedCountry 代表国家
	FieldRepresentedCountry = "represented_country"

	// FieldRepresentedCountryType 代表国家类型
	FieldRepresentedCountryType = "represented_country_type"

	// FieldIsAnonymousProxy 是否匿名代理
	FieldIsAnonymousProxy = "is_anonymous_proxy"

	// FieldIsSatelliteProvider 是否卫星提供商
	FieldIsSatelliteProvider = "is_satellite_provider"

	// ASN Fields

	// FieldAutonomousSystemNumber 自治系统号
	FieldAutonomousSystemNumber = "autonomous_system_number"

	// FieldAutonomousSystemOrganization 自治系统组织
	FieldAutonomousSystemOrganization = "autonomous_system_organization"
)

// CommonFieldsAlias 公共字段到数据库字段映射
var CommonFieldsAlias = map[string]string{
	model.Country:   FieldCountry,
	model.City:      FieldCity,
	model.Continent: FieldContinent,
	model.Province:  FieldSubdivisions,
	model.UTCOffset: FieldTimeZone,
	model.Latitude:  FieldLatitude,
	model.Longitude: FieldLongitude,
	model.ASN:       FieldAutonomousSystemNumber,
}
