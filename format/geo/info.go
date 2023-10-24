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

package geo

import (
	"strconv"
	"strings"

	"github.com/sjzar/ips/format/geo/data"
)

// Info represents the geolocation information.
type Info struct {
	GeoNameID         int               `json:"geoname_id"`           // GeoNameID is the ID of the record in the GeoNames database.
	Names             map[string]string `json:"names"`                // Names is map of locale codes to the name in that locale.
	Code              string            `json:"code"`                 // Code is the Code for Continent.
	IsoCode           string            `json:"iso_code"`             // IsoCode is the ISO Code for the country or subdivision.
	IsInEuropeanUnion bool              `json:"is_in_european_union"` // IsInEuropeanUnion is true if the country is a member state of the European Union.
	ConnectionID      int               `json:"-"`                    // ConnectionID is the GeoNameID for the Connection Type.
	CountryID         int               `json:"-"`                    // CountryID is the GeoNameID for the Country.
	SubdivisionIDs    []int             `json:"-"`                    // SubdivisionIDs is an array of GeoNameIDs for the Subdivisions(regions).
}

// Name returns the name of the geolocation info in the specified language.
func (g *Info) Name(lang string) string {
	if len(g.Names) == 0 {
		return ""
	}

	if v, ok := g.Names[lang]; ok {
		return v
	}
	if v, ok := g.Names[Language]; ok {
		return v
	}
	if v, ok := g.Names[LangEnglish]; ok {
		return v
	}
	for _, v := range g.Names {
		return v
	}

	// impossible
	return ""
}

// Map converts the geolocation info into a map structure.
func (g *Info) Map(selectLanguages string) map[string]interface{} {
	ret := make(map[string]interface{})
	ret["geoname_id"] = g.GeoNameID
	if selectLanguages != "-" {
		langMap := make(map[string]bool)
		for _, lang := range strings.Split(selectLanguages, ",") {
			if len(lang) == 0 {
				continue
			}
			langMap[lang] = true
		}
		names := make(map[string]interface{})
		for k, v := range g.Names {
			if len(langMap) != 0 && !langMap[k] {
				continue
			}
			names[k] = v
		}
		ret["names"] = names
	}
	if len(g.Code) != 0 {
		ret["code"] = g.Code
	}
	if len(g.IsoCode) != 0 {
		ret["iso_code"] = g.IsoCode
	}
	if g.IsInEuropeanUnion {
		ret["is_in_european_union"] = true
	}
	return ret
}

// ParseInfoFromMMDB extracts geolocation info from a MMDB map.
func ParseInfoFromMMDB(m map[string]interface{}, disableExtraData bool) (*Info, bool) {
	ret := &Info{}
	for key, value := range m {
		switch key {
		case "geoname_id":
			if v, ok := value.(int); ok {
				ret.GeoNameID = v
			}
			if !disableExtraData {
				if extra, ok := GetInfoByID(ret.GeoNameID); ok {
					return extra, true
				}
			}
		case "names":
			if v, ok := value.(map[string]interface{}); ok {
				ret.Names = make(map[string]string)
				for key2, value2 := range v {
					if v2, ok := value2.(string); ok {
						ret.Names[key2] = v2
					}
				}
			}
		case "code":
			if v, ok := value.(string); ok {
				ret.Code = v
			}
		case "iso_code":
			if v, ok := value.(string); ok {
				ret.IsoCode = v
			}
		case "is_in_european_union":
			if v, ok := value.(bool); ok {
				ret.IsInEuropeanUnion = v
			}
		}
	}

	if len(ret.Names) == 0 {
		return nil, false
	}

	return ret, true
}

// ParseGeoInfo parses a geolocation info from a formatted string.
// Format: GeoNameID, Names, Code, ISO-Code, IsInEuropeanUnion, ConnectionID, CountryID, SubdivisionIDs
func ParseGeoInfo(s string) (*Info, bool) {
	split := strings.Split(s, "\t")
	ret := &Info{}
	for i, v := range split {
		switch i {
		case 0:
			ret.GeoNameID, _ = strconv.Atoi(v)
		case 1:
			ret.Names = make(map[string]string)
			split2 := strings.Split(v, "|")
			for _, v2 := range split2 {
				if len(v2) == 0 {
					continue
				}
				split3 := strings.SplitN(v2, ":", 2)
				if len(split3) == 2 {
					ret.Names[split3[0]] = split3[1]
				}
			}
		case 2:
			ret.Code = v
		case 3:
			ret.IsoCode = v
		case 4:
			ret.IsInEuropeanUnion = v == "1"
		case 5:
			ret.ConnectionID, _ = strconv.Atoi(v)
		case 6:
			ret.CountryID, _ = strconv.Atoi(v)
		case 7:
			split2 := strings.Split(v, "|")
			for _, v2 := range split2 {
				if len(v2) == 0 {
					continue
				}
				v3, _ := strconv.Atoi(v2)
				ret.SubdivisionIDs = append(ret.SubdivisionIDs, v3)
			}
		}
	}
	if len(ret.Names) == 0 {
		return nil, false
	}
	return ret, true
}

// GetInfoByID retrieves geolocation info by its GeoNameID.
func GetInfoByID(geoNameID int) (*Info, bool) {
	if IDInfos == nil {
		IDInfos = LoadData("", data.Continent, data.Country, data.Region, data.City)
	}

	str, ok := IDInfos[strconv.Itoa(geoNameID)]
	if !ok {
		return nil, false
	}
	return ParseGeoInfo(str)
}

// GetInfoByName retrieves geolocation info by its name.
func GetInfoByName(field, name string) (*Info, bool) {
	nameInfos := GetNameInfos(field, Language)

	str, ok := nameInfos[name]
	if !ok {
		return nil, false
	}
	return ParseGeoInfo(str)
}
