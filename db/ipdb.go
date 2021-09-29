package db

import (
	"log"
	"net"

	"github.com/sjzar/ips/db/ipdb"
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
	IPDBFieldCountryName    = "country_name"
	IPDBFieldRegionName     = "region_name"
	IPDBFieldCityName       = "city_name"
	IPDBFieldISPDomain      = "isp_domain"
	IPDBFieldContinentCode  = "continent_code"
	IPDBFieldUTCOffset      = "utc_offset"
	IPDBFieldLatitude       = "latitude"
	IPDBFieldLongitude      = "longitude"
	IPDBFieldChinaAdminCode = "china_admin_code"
	IPDBFieldOwnerDomain    = "owner_domain"
	IPDBFieldTimezone       = "timezone"
	IPDBFieldIddCode        = "idd_code"
	IPDBFieldCountryCode    = "country_code"
	IPDBFieldIDC            = "idc"
	IPDBFieldBaseStation    = "base_station"
	IPDBFieldCountryCode3   = "country_code3"
	IPDBFieldEuropeanUnion  = "european_union"
	IPDBFieldCurrencyCode   = "currency_code"
	IPDBFieldCurrencyName   = "currency_name"
	IPDBFieldAnycast        = "anycast"
)

// IPDBFieldsMap IPDB字段映射表
var IPDBFieldsMap = map[string]string{
	model.Country:          IPDBFieldCountryName,
	model.Province:         IPDBFieldRegionName,
	model.City:             IPDBFieldCityName,
	model.ISP:              IPDBFieldISPDomain,
	model.Continent:        IPDBFieldContinentCode,
	model.UTCOffset:        IPDBFieldUTCOffset,
	model.Latitude:         IPDBFieldLatitude,
	model.Longitude:        IPDBFieldLongitude,
	model.ChinaAdminCode:   IPDBFieldChinaAdminCode,
	IPDBFieldOwnerDomain:   IPDBFieldOwnerDomain,
	IPDBFieldTimezone:      IPDBFieldTimezone,
	IPDBFieldIddCode:       IPDBFieldIddCode,
	IPDBFieldCountryCode:   IPDBFieldCountryCode,
	IPDBFieldIDC:           IPDBFieldIDC,
	IPDBFieldBaseStation:   IPDBFieldBaseStation,
	IPDBFieldCountryCode3:  IPDBFieldCountryCode3,
	IPDBFieldEuropeanUnion: IPDBFieldEuropeanUnion,
	IPDBFieldCurrencyCode:  IPDBFieldCurrencyCode,
	IPDBFieldCurrencyName:  IPDBFieldCurrencyName,
	IPDBFieldAnycast:       IPDBFieldAnycast,
}

var IPDBFullFields = []string{
	model.Country, model.Province, model.City, model.ISP,
	model.Continent, model.UTCOffset, model.Latitude, model.Longitude, model.ChinaAdminCode,
	IPDBFieldOwnerDomain, IPDBFieldTimezone, IPDBFieldIddCode, IPDBFieldCountryCode,
	IPDBFieldIDC, IPDBFieldBaseStation, IPDBFieldCountryCode3, IPDBFieldEuropeanUnion,
	IPDBFieldCurrencyCode, IPDBFieldCurrencyName, IPDBFieldAnycast,
}

func IPDBFieldsMapping(fields []string) []string {
	ret := make([]string, len(fields))
	for i := range fields {
		ret[i] = IPDBFieldsMap[fields[i]]
	}
	return ret
}

type IPDB struct {
	db *ipdb.City
}

func NewIPDB(file string) *IPDB {
	city, err := ipdb.NewCity(file)
	if err != nil {
		log.Fatal(err)
	}

	return &IPDB{
		db: city,
	}
}

func (i *IPDB) LookupNetwork(ip net.IP) (*net.IPNet, map[string]string, error) {
	cityInfo, ipNet, err := i.db.FindInfo(ip.String(), "CN")
	if err != nil {
		return nil, nil, err
	}

	return ipNet, i.format(cityInfo), nil
}
func (i *IPDB) format(cityInfo *ipdb.CityInfo) map[string]string {
	data := map[string]string{
		model.Country:          cityInfo.CountryName,
		model.Province:         cityInfo.RegionName,
		model.City:             cityInfo.CityName,
		model.ISP:              cityInfo.IspDomain,
		model.Continent:        cityInfo.ContinentCode,
		model.UTCOffset:        cityInfo.UtcOffset,
		model.Latitude:         cityInfo.Latitude,
		model.Longitude:        cityInfo.Longitude,
		model.ChinaAdminCode:   cityInfo.ChinaAdminCode,
		IPDBFieldOwnerDomain:   cityInfo.OwnerDomain,
		IPDBFieldTimezone:      cityInfo.Timezone,
		IPDBFieldIddCode:       cityInfo.IddCode,
		IPDBFieldCountryCode:   cityInfo.CountryCode,
		IPDBFieldIDC:           cityInfo.IDC,
		IPDBFieldBaseStation:   cityInfo.BaseStation,
		IPDBFieldCountryCode3:  cityInfo.CountryCode3,
		IPDBFieldEuropeanUnion: cityInfo.EuropeanUnion,
		IPDBFieldCurrencyCode:  cityInfo.CurrencyCode,
		IPDBFieldCurrencyName:  cityInfo.CurrencyName,
		IPDBFieldAnycast:       cityInfo.Anycast,
	}

	return data
}
