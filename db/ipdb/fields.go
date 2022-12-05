package ipdb

import (
	"github.com/sjzar/ips/model"
)

// "country_name": "中国",
// "region_name": "浙江",
// "city_name": "",
// "isp_domain": "电信",
// "continent_code": "AP",
// "utc_offset": "UTC+8",
// "latitude": "29.19083",
// "longitude": "120.083656",
// "china_admin_code": "330000",
// "owner_domain": "",
// "timezone": "Asia/Shanghai",
// "idd_code": "86",
// "country_code": "CN",

const (
	FieldCountryName    = "country_name"
	FieldRegionName     = "region_name"
	FieldCityName       = "city_name"
	FieldISPDomain      = "isp_domain"
	FieldContinentCode  = "continent_code"
	FieldUTCOffset      = "utc_offset"
	FieldLatitude       = "latitude"
	FieldLongitude      = "longitude"
	FieldChinaAdminCode = "china_admin_code"
	FieldOwnerDomain    = "owner_domain"
	FieldTimezone       = "timezone"
	FieldIddCode        = "idd_code"
	FieldCountryCode    = "country_code"
	FieldIDC            = "idc"
	FieldBaseStation    = "base_station"
	FieldCountryCode3   = "country_code3"
	FieldEuropeanUnion  = "european_union"
	FieldCurrencyCode   = "currency_code"
	FieldCurrencyName   = "currency_name"
	FieldAnycast        = "anycast"
)

// FullFields 全字段列表
var FullFields = []string{
	FieldCountryName,
	FieldRegionName,
	FieldCityName,
	FieldISPDomain,
	FieldContinentCode,
	FieldUTCOffset,
	FieldLatitude,
	FieldLongitude,
	FieldChinaAdminCode,
	FieldOwnerDomain,
	FieldTimezone,
	FieldIddCode,
	FieldCountryCode,
	FieldIDC,
	FieldBaseStation,
	FieldCountryCode3,
	FieldEuropeanUnion,
	FieldCurrencyCode,
	FieldCurrencyName,
	FieldAnycast,
}

// CommonFieldsMap 公共字段映射
var CommonFieldsMap = map[string]string{
	model.Country:        FieldCountryName,
	model.Province:       FieldRegionName,
	model.City:           FieldCityName,
	model.ISP:            FieldISPDomain,
	model.Continent:      FieldContinentCode,
	model.UTCOffset:      FieldUTCOffset,
	model.Latitude:       FieldLatitude,
	model.Longitude:      FieldLongitude,
	model.ChinaAdminCode: FieldChinaAdminCode,
}

// FieldsFormat 字段格式化，并补充公共字段
func FieldsFormat(data map[string]string) map[string]string {

	// Fill Common Fields
	for k, v := range CommonFieldsMap {
		data[k] = data[v]
	}

	return data
}

// FieldsReplace 字段替换
func FieldsReplace(fields []string) []string {
	for i := range fields {
		if v, ok := CommonFieldsMap[fields[i]]; ok {
			fields[i] = v
		}
	}
	return fields
}
