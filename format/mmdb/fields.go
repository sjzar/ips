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
	"strconv"

	"github.com/maxmind/mmdbwriter/mmdbtype"

	"github.com/sjzar/ips/pkg/model"
)

var Lang = "zh-CN"

const (

	// FieldCity 城市
	FieldCity = "city"

	// FieldContinent 大洲
	FieldContinent = "continent"

	// FieldContinentCode 大洲代码
	FieldContinentCode = "continent_code"

	// FieldCountry 国家
	FieldCountry = "country"

	// FieldCountryISOCode 国家ISO代码
	FieldCountryISOCode = "country_iso_code"

	// FieldCountryIsInEuropeanUnion 国家是否在欧盟
	FieldCountryIsInEuropeanUnion = "country_is_in_european_union"

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

	// FieldRegisteredCountryISOCode 注册国家ISO代码
	FieldRegisteredCountryISOCode = "registered_country_iso_code"

	// FieldRegisteredCountryIsInEuropeanUnion 注册国家是否在欧盟
	FieldRegisteredCountryIsInEuropeanUnion = "registered_country_is_in_european_union"

	// FieldRepresentedCountry 代表国家
	FieldRepresentedCountry = "represented_country"

	// FieldRepresentedCountryISOCode 代表国家ISO代码
	FieldRepresentedCountryISOCode = "represented_country_iso_code"

	// FieldRepresentedCountryIsInEuropeanUnion 代表国家是否在欧盟
	FieldRepresentedCountryIsInEuropeanUnion = "represented_country_is_in_european_union"

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

func (c *City) Format() map[string]string {
	data := make(map[string]string)
	if c.City.Names != nil {
		data[FieldCity] = c.City.Names[Lang]
	}

	if c.Continent.Names != nil {
		data[FieldContinent] = c.Continent.Names[Lang]
	}
	data[FieldContinentCode] = c.Continent.Code

	if c.Country.Names != nil {
		data[FieldCountry] = c.Country.Names[Lang]
	}
	data[FieldCountryISOCode] = c.Country.IsoCode
	data[FieldCountryIsInEuropeanUnion] = strconv.FormatBool(c.Country.IsInEuropeanUnion)

	data[FieldAccuracyRadius] = strconv.Itoa(int(c.Location.AccuracyRadius))
	data[FieldLatitude] = strconv.FormatFloat(c.Location.Latitude, 'f', -1, 64)
	data[FieldLongitude] = strconv.FormatFloat(c.Location.Longitude, 'f', -1, 64)
	data[FieldMetroCode] = strconv.Itoa(int(c.Location.MetroCode))
	data[FieldTimeZone] = c.Location.TimeZone
	data[FieldPostalCode] = c.Postal.Code

	if c.RegisteredCountry.Names != nil {
		data[FieldRegisteredCountry] = c.RegisteredCountry.Names[Lang]
	}
	data[FieldRegisteredCountryISOCode] = c.RegisteredCountry.IsoCode
	data[FieldRegisteredCountryIsInEuropeanUnion] = strconv.FormatBool(c.RegisteredCountry.IsInEuropeanUnion)

	// if c.RepresentedCountry.Names != nil {
	//	data[FieldRepresentedCountry] = c.RepresentedCountry.Names[Lang]
	// }
	// data[FieldRepresentedCountryISOCode] = c.RepresentedCountry.IsoCode
	// data[FieldRepresentedCountryIsInEuropeanUnion] = strconv.FormatBool(c.RepresentedCountry.IsInEuropeanUnion)
	// data[FieldRepresentedCountryType] = c.RepresentedCountry.Type

	data[FieldIsAnonymousProxy] = strconv.FormatBool(c.Traits.IsAnonymousProxy)
	data[FieldIsSatelliteProvider] = strconv.FormatBool(c.Traits.IsSatelliteProvider)

	return data
}

func (a *ASN) Format() map[string]string {
	data := make(map[string]string)
	data[FieldAutonomousSystemNumber] = strconv.FormatUint(uint64(a.AutonomousSystemNumber), 10)
	data[FieldAutonomousSystemOrganization] = a.AutonomousSystemOrganization

	return data
}

// ConvertMap converts fields and values to mmdbtype.Map
// FIXME temporary solution :)
func ConvertMap(fields, values []string) mmdbtype.DataType {
	ret := mmdbtype.Map{}
	for i := range fields {
		if len(values[i]) == 0 {
			continue
		}
		switch fields[i] {
		case FieldCity, FieldContinent, FieldCountry, FieldRegisteredCountry:
			if dataMap, ok := getMap(ret, fields[i]); ok {
				dataMap.(mmdbtype.Map)["names"] = mmdbtype.Map{
					mmdbtype.String(Lang): mmdbtype.String(values[i]),
				}
			}
		case FieldContinentCode:
			if dataMap, ok := getMap(ret, FieldContinent); ok {
				dataMap.(mmdbtype.Map)["code"] = mmdbtype.String(values[i])
			}
		case FieldCountryISOCode:
			if dataMap, ok := getMap(ret, FieldCountry); ok {
				dataMap.(mmdbtype.Map)["iso_code"] = mmdbtype.String(values[i])
			}
		case FieldCountryIsInEuropeanUnion:
			value, err := strconv.ParseBool(values[i])
			if err != nil {
				continue
			}
			if dataMap, ok := getMap(ret, FieldCountry); ok {
				dataMap.(mmdbtype.Map)["is_in_european_union"] = mmdbtype.Bool(value)
			}
		case FieldAccuracyRadius:
			value, err := strconv.ParseUint(values[i], 10, 16)
			if err != nil {
				continue
			}
			if dataMap, ok := getMap(ret, "location"); ok {
				dataMap.(mmdbtype.Map)[FieldAccuracyRadius] = mmdbtype.Uint16(value)
			}
		case FieldLatitude:
			value, err := strconv.ParseFloat(values[i], 64)
			if err != nil {
				continue
			}
			if dataMap, ok := getMap(ret, "location"); ok {
				dataMap.(mmdbtype.Map)[FieldLatitude] = mmdbtype.Float64(value)
			}
		case FieldLongitude:
			value, err := strconv.ParseFloat(values[i], 64)
			if err != nil {
				continue
			}
			if dataMap, ok := getMap(ret, "location"); ok {
				dataMap.(mmdbtype.Map)[FieldLongitude] = mmdbtype.Float64(value)
			}
		case FieldMetroCode:
			value, err := strconv.ParseUint(values[i], 10, 64)
			if err != nil {
				continue
			}
			if dataMap, ok := getMap(ret, "location"); ok {
				dataMap.(mmdbtype.Map)[FieldMetroCode] = mmdbtype.Uint64(value)
			}
		case FieldTimeZone:
			if dataMap, ok := getMap(ret, "location"); ok {
				dataMap.(mmdbtype.Map)[FieldTimeZone] = mmdbtype.String(values[i])
			}
		case FieldPostalCode:
			ret[mmdbtype.String("postal")] = mmdbtype.Map{
				"code": mmdbtype.String(values[i]),
			}
		case FieldRegisteredCountryISOCode:
			if dataMap, ok := getMap(ret, FieldRegisteredCountry); ok {
				dataMap.(mmdbtype.Map)["iso_code"] = mmdbtype.String(values[i])
			}
		case FieldRegisteredCountryIsInEuropeanUnion:
			value, err := strconv.ParseBool(values[i])
			if err != nil {
				continue
			}
			if dataMap, ok := getMap(ret, FieldRegisteredCountry); ok {
				dataMap.(mmdbtype.Map)["is_in_european_union"] = mmdbtype.Bool(value)
			}
		case FieldIsAnonymousProxy:
			value, err := strconv.ParseBool(values[i])
			if err != nil {
				continue
			}
			if dataMap, ok := getMap(ret, "traits"); ok {
				dataMap.(mmdbtype.Map)[FieldIsAnonymousProxy] = mmdbtype.Bool(value)
			}
		case FieldIsSatelliteProvider:
			value, err := strconv.ParseBool(values[i])
			if err != nil {
				continue
			}
			if dataMap, ok := getMap(ret, "traits"); ok {
				dataMap.(mmdbtype.Map)[FieldIsSatelliteProvider] = mmdbtype.Bool(value)
			}
		case FieldAutonomousSystemNumber:
			value, err := strconv.ParseUint(values[i], 10, 64)
			if err != nil {
				continue
			}
			ret[FieldAutonomousSystemNumber] = mmdbtype.Uint64(value)
		case FieldAutonomousSystemOrganization:
			ret[FieldAutonomousSystemOrganization] = mmdbtype.String(values[i])
		default:
			ret[mmdbtype.String(fields[i])] = mmdbtype.String(values[i])
		}
	}

	return ret
}

// getMap retrieves a map from the provided mmdbtype.DataType.
func getMap(data mmdbtype.DataType, key string) (mmdbtype.DataType, bool) {
	if m, ok := data.(mmdbtype.Map); ok {
		if v, ok := m[mmdbtype.String(key)]; ok {
			return v, true
		}
		newMap := mmdbtype.Map{}
		m[mmdbtype.String(key)] = newMap
		return newMap, true
	}
	return nil, false
}

// CityFullFields 城市数据库全字段列表
var CityFullFields = []string{
	FieldCity,
	FieldContinent,
	FieldContinentCode,
	FieldCountry,
	FieldCountryISOCode,
	FieldCountryIsInEuropeanUnion,
	FieldAccuracyRadius,
	FieldLatitude,
	FieldLongitude,
	FieldMetroCode,
	FieldTimeZone,
	FieldPostalCode,
	FieldRegisteredCountry,
	FieldRegisteredCountryISOCode,
	FieldRegisteredCountryIsInEuropeanUnion,
	// FieldRepresentedCountry,
	// FieldRepresentedCountryISOCode,
	// FieldRepresentedCountryIsInEuropeanUnion,
	// FieldRepresentedCountryType,
	FieldIsAnonymousProxy,
	FieldIsSatelliteProvider,
}

// ASNFullFields ASN数据库全字段列表
var ASNFullFields = []string{
	FieldAutonomousSystemNumber,
	FieldAutonomousSystemOrganization,
}

// CommonFieldsAlias 公共字段到数据库字段映射
var CommonFieldsAlias = map[string]string{
	model.Country:   FieldCountry,
	model.City:      FieldCity,
	model.Continent: FieldContinent,
	model.UTCOffset: FieldTimeZone,
	model.Latitude:  FieldLatitude,
	model.Longitude: FieldLongitude,
	model.ASN:       FieldAutonomousSystemNumber,
}
