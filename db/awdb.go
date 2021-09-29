package db

import (
	"log"
	"net"

	"github.com/dilfish/awdb-golang/awdb-golang"

	"github.com/sjzar/ips/iprange"
	"github.com/sjzar/ips/mapper"
	"github.com/sjzar/ips/model"
)

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

// awdb 参数
const (
	AWDBFieldCountry   = "country"
	AWDBFieldProvince  = "province"
	AWDBFieldCity      = "city"
	AWDBFieldISP       = "isp"
	AWDBFieldContinent = "continent"
	AWDBFieldTimeZone  = "timezone"
	AWDBFieldLatwgs    = "latwgs"
	AWDBFieldLngwgs    = "lngwgs"
	AWDBFieldAdcode    = "adcode"
	AWDBFieldAccuracy  = "accuracy"
	AWDBFieldAreaCode  = "areacode"
	AWDBFieldASNumber  = "asnumber"
	AWDBFieldOwner     = "owner"
	AWDBFieldRadius    = "radius"
	AWDBFieldSource    = "source"
	AWDBFieldZipCode   = "zipcode"
)

// AWDBFieldsMap AWDB字段映射表
var AWDBFieldsMap = map[string]string{
	model.Country:        AWDBFieldCountry,
	model.Province:       AWDBFieldProvince,
	model.City:           AWDBFieldCity,
	model.ISP:            AWDBFieldISP,
	model.Continent:      AWDBFieldContinent,
	model.UTCOffset:      AWDBFieldTimeZone,
	model.Latitude:       AWDBFieldLatwgs,
	model.Longitude:      AWDBFieldLngwgs,
	model.ChinaAdminCode: AWDBFieldAdcode,
	AWDBFieldAccuracy:    AWDBFieldAccuracy,
	AWDBFieldAreaCode:    AWDBFieldAreaCode,
	AWDBFieldASNumber:    AWDBFieldASNumber,
	AWDBFieldOwner:       AWDBFieldOwner,
	AWDBFieldRadius:      AWDBFieldRadius,
	AWDBFieldSource:      AWDBFieldSource,
	AWDBFieldZipCode:     AWDBFieldZipCode,
}

var AWDBFullFields = []string{
	model.Country, model.Province, model.City, model.ISP,
	model.Continent, model.UTCOffset, model.Latitude, model.Longitude, model.ChinaAdminCode,
	AWDBFieldAccuracy, AWDBFieldAreaCode, AWDBFieldASNumber, AWDBFieldOwner,
	AWDBFieldRadius, AWDBFieldSource, AWDBFieldZipCode,
}

type AWDB struct {
	db       *awdb.Reader
	ipLineDB *awdb.Reader

	// extra mapper
	ASNMapper mapper.Mapper
}

func NewAWDB(file, iplineFile string) *AWDB {
	db, err := awdb.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	var ipLineDB *awdb.Reader
	if len(iplineFile) != 0 {
		ipLineDB, err = awdb.Open(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &AWDB{
		db:       db,
		ipLineDB: ipLineDB,
	}
}

func (a *AWDB) Close() {
	if a.db != nil {
		a.db.Close()
	}
	if a.ipLineDB != nil {
		a.ipLineDB.Close()
	}
}

func (a *AWDB) LookupIPLine(ip net.IP) (*net.IPNet, string) {
	if a.ipLineDB == nil {
		return nil, ""
	}
	var record interface{}
	ipNet, ok, err := a.ipLineDB.LookupNetwork(ip, &record)
	if !ok || err != nil {
		return nil, ""
	}
	_record := record.(map[string]interface{})
	if _line, ok := _record["line"]; ok {
		line := string(_line.([]byte))
		return ipNet, line
	}
	return nil, ""
}

func (a *AWDB) LookupNetwork(ip net.IP) (*net.IPNet, map[string]string, error) {
	var record interface{}
	ipNet, ok, err := a.db.LookupNetwork(ip, &record)
	if !ok || err != nil {
		return nil, nil, err
	}
	ipNetLine, line := a.LookupIPLine(ip)
	if len(line) != 0 {
		if !iprange.IPNetMaskLess(ipNet, ipNetLine) {
			ipNet = ipNetLine
		}
	}

	return ipNet, a.format(record.(map[string]interface{}), line), nil
}

func (a *AWDB) format(record map[string]interface{}, line string) map[string]string {
	data := make(map[string]string)

	for k, v := range AWDBFieldsMap {
		_v, ok := record[v]
		if !ok {
			continue
		}
		data[k] = string(_v.([]byte))

		// line Mapping && ASN ISP Mapping
		if k == model.ISP {
			if len(line) != 0 {
				data[k] = line
			} else if asn, ok := record[AWDBFieldASNumber]; ok && a.ASNMapper != nil {
				if isp, ok := a.ASNMapper.Mapping("", string(asn.([]byte))); ok {
					data[k] = isp
				}
			}
		}
	}

	return data
}
