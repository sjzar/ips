package mmdb

import (
	"strconv"

	"github.com/sjzar/ips/model"
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
)

// The City struct corresponds to the data in the GeoIP2/GeoLite2 City
// databases.
// Copy From https://github.com/oschwald/geoip2-golang/blob/ab86d30fa7a001e932981dab47fc57fb983048d5/reader.go#L87
type City struct {
	City struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Continent struct {
		Code      string            `maxminddb:"code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"continent"`
	Country struct {
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
		IsoCode           string            `maxminddb:"iso_code"`
		Names             map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Location struct {
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		MetroCode      uint    `maxminddb:"metro_code"`
		TimeZone       string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	RegisteredCountry struct {
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
		IsoCode           string            `maxminddb:"iso_code"`
		Names             map[string]string `maxminddb:"names"`
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
		IsoCode           string            `maxminddb:"iso_code"`
		Names             map[string]string `maxminddb:"names"`
		Type              string            `maxminddb:"type"`
	} `maxminddb:"represented_country"`
	Subdivisions []struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
}

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

	//if c.RepresentedCountry.Names != nil {
	//	data[FieldRepresentedCountry] = c.RepresentedCountry.Names[Lang]
	//}
	//data[FieldRepresentedCountryISOCode] = c.RepresentedCountry.IsoCode
	//data[FieldRepresentedCountryIsInEuropeanUnion] = strconv.FormatBool(c.RepresentedCountry.IsInEuropeanUnion)
	//data[FieldRepresentedCountryType] = c.RepresentedCountry.Type

	data[FieldIsAnonymousProxy] = strconv.FormatBool(c.Traits.IsAnonymousProxy)
	data[FieldIsSatelliteProvider] = strconv.FormatBool(c.Traits.IsSatelliteProvider)

	return FieldsFormat(data)
}

// FullFields 全字段列表
var FullFields = []string{
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
	//FieldRepresentedCountry,
	//FieldRepresentedCountryISOCode,
	//FieldRepresentedCountryIsInEuropeanUnion,
	//FieldRepresentedCountryType,
	FieldIsAnonymousProxy,
	FieldIsSatelliteProvider,
}

// CommonFieldsMap 公共字段映射
var CommonFieldsMap = map[string]string{
	model.Country:   FieldCountry,
	model.City:      FieldCity,
	model.Continent: FieldContinent,
	model.UTCOffset: FieldTimeZone,
	model.Latitude:  FieldLatitude,
	model.Longitude: FieldLongitude,
}

// FieldsFormat 字段格式化，并补充公共字段
func FieldsFormat(record map[string]string) map[string]string {
	data := make(map[string]string)

	for k, v := range record {
		data[k] = v
	}

	// Fill Common Fields
	for k, v := range CommonFieldsMap {
		data[k] = data[v]
	}

	return data
}
