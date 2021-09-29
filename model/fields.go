package model

const (
	// Country 国家
	Country = "country"

	// Province 省份
	Province = "province"

	// City 城市
	City = "city"

	// ISP 运营商
	ISP = "isp"

	// Continent 大区
	Continent = "continent"

	// UTCOffset UTC偏移值
	UTCOffset = "utcOffset"

	// Latitude 纬度
	Latitude = "latitude"

	// Longitude 经度
	Longitude = "longitude"

	// ChinaAdminCode 中国行政区划代码
	ChinaAdminCode = "chinaAdminCode"

	// Placeholder 占位符
	Placeholder = "-"
)

var FullFields = []string{
	Country, Province, City, ISP,
	Continent, UTCOffset, Latitude, Longitude,
	ChinaAdminCode,
}
