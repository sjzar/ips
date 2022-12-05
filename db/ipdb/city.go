package ipdb

// Copy From https://github.com/ipipdotnet/ipdb-go
// modify by sjzar
// modify: Find / FindMap return ipNet

import (
	"encoding/json"
	"net"
	"os"
	"reflect"
	"time"
)

// CityInfo is City Database Content
type CityInfo struct {
	CountryName    string `json:"country_name"`
	RegionName     string `json:"region_name"`
	CityName       string `json:"city_name"`
	OwnerDomain    string `json:"owner_domain"`
	IspDomain      string `json:"isp_domain"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	Timezone       string `json:"timezone"`
	UtcOffset      string `json:"utc_offset"`
	ChinaAdminCode string `json:"china_admin_code"`
	IddCode        string `json:"idd_code"`
	CountryCode    string `json:"country_code"`
	ContinentCode  string `json:"continent_code"`
	IDC            string `json:"idc"`
	BaseStation    string `json:"base_station"`
	CountryCode3   string `json:"country_code3"`
	EuropeanUnion  string `json:"european_union"`
	CurrencyCode   string `json:"currency_code"`
	CurrencyName   string `json:"currency_name"`
	Anycast        string `json:"anycast"`

	Line string `json:"line"`

	DistrictInfo DistrictInfo `json:"district_info"`

	Route   string    `json:"route"`
	ASN     string    `json:"asn"`
	ASNInfo []ASNInfo `json:"asn_info"`

	AreaCode string `json:"area_code"`

	UsageType string `json:"usage_type"`
}

type ASNInfo struct {
	ASN      int    `json:"asn"`
	Registry string `json:"reg"`
	Country  string `json:"cc"`
	Net      string `json:"net"`
	Org      string `json:"org"`
	Type     string `json:"type"`
	Domain   string `json:"domain"`
}

type DistrictInfo struct {
	CountryName    string `json:"country_name"`
	RegionName     string `json:"region_name"`
	CityName       string `json:"city_name"`
	DistrictName   string `json:"district_name"`
	ChinaAdminCode string `json:"china_admin_code"`
	CoveringRadius string `json:"covering_radius"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
}

// City struct
type City struct {
	reader *reader
}

// NewCity initialize
func NewCity(name string) (*City, error) {

	r, e := newReader(name, &CityInfo{})
	if e != nil {
		return nil, e
	}

	return &City{
		reader: r,
	}, nil
}

// Reload the database
func (db *City) Reload(name string) error {

	_, err := os.Stat(name)
	if err != nil {
		return err
	}

	reader, err := newReader(name, &CityInfo{})
	if err != nil {
		return err
	}

	db.reader = reader

	return nil
}

// Find query with addr
func (db *City) Find(addr, language string) ([]string, *net.IPNet, error) {
	return db.reader.find1(addr, language)
}

// FindMap query with addr
func (db *City) FindMap(addr, language string) (map[string]string, *net.IPNet, error) {

	data, ipNet, err := db.reader.find1(addr, language)
	if err != nil {
		return nil, nil, err
	}
	info := make(map[string]string, len(db.reader.meta.Fields))
	for k, v := range data {
		info[db.reader.meta.Fields[k]] = v
	}

	return info, ipNet, nil
}

// FindInfo query with addr
func (db *City) FindInfo(addr, language string) (*CityInfo, *net.IPNet, error) {

	data, ipNet, err := db.reader.FindMap(addr, language)
	if err != nil {
		return nil, nil, err
	}

	var asnInfoList []ASNInfo
	var asnInfoType = reflect.TypeOf(asnInfoList)

	var districtInfo DistrictInfo
	var districtInfoType = reflect.TypeOf(districtInfo)

	info := &CityInfo{}

	for k, v := range data {
		sv := reflect.ValueOf(info).Elem()
		sfv := sv.FieldByName(db.reader.refType[k])

		if !sfv.IsValid() {
			continue
		}
		if !sfv.CanSet() {
			continue
		}

		sft := sfv.Type()
		fv := reflect.ValueOf(v)
		if sft == fv.Type() {
			sfv.Set(fv)
		} else if sft == asnInfoType {
			err = json.Unmarshal([]byte(v), &asnInfoList)
			if err == nil {
				sfv.Set(reflect.ValueOf(asnInfoList))
			}
		} else if sft == districtInfoType {
			err = json.Unmarshal([]byte(v), &districtInfo)
			if err == nil {
				sfv.Set(reflect.ValueOf(districtInfo))
			}
		}
	}

	return info, ipNet, nil
}

// IsIPv4 whether support ipv4
func (db *City) IsIPv4() bool {
	return db.reader.IsIPv4Support()
}

// IsIPv6 whether support ipv6
func (db *City) IsIPv6() bool {
	return db.reader.IsIPv6Support()
}

// Languages return support languages
func (db *City) Languages() []string {
	return db.reader.Languages()
}

// Fields return support fields
func (db *City) Fields() []string {
	return db.reader.meta.Fields
}

// BuildTime return database build Time
func (db *City) BuildTime() time.Time {
	return db.reader.Build()
}
