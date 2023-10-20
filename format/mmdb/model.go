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

// Copy From https://github.com/oschwald/geoip2-golang/blob/ab86d30fa7a001e932981dab47fc57fb983048d5/reader.go

// The Enterprise struct corresponds to the data in the GeoIP2 Enterprise
// database.
type Enterprise struct {
	City struct {
		Confidence uint8             `maxminddb:"confidence"`
		GeoNameID  uint              `maxminddb:"geoname_id"`
		Names      map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Continent struct {
		Code      string            `maxminddb:"code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"continent"`
	Country struct {
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsoCode           string            `maxminddb:"iso_code"`
		Names             map[string]string `maxminddb:"names"`
		Confidence        uint8             `maxminddb:"confidence"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"country"`
	Location struct {
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		MetroCode      uint    `maxminddb:"metro_code"`
		TimeZone       string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
	Postal struct {
		Code       string `maxminddb:"code"`
		Confidence uint8  `maxminddb:"confidence"`
	} `maxminddb:"postal"`
	RegisteredCountry struct {
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsoCode           string            `maxminddb:"iso_code"`
		Names             map[string]string `maxminddb:"names"`
		Confidence        uint8             `maxminddb:"confidence"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
		IsoCode           string            `maxminddb:"iso_code"`
		Names             map[string]string `maxminddb:"names"`
		Type              string            `maxminddb:"type"`
	} `maxminddb:"represented_country"`
	Subdivisions []struct {
		Confidence uint8             `maxminddb:"confidence"`
		GeoNameID  uint              `maxminddb:"geoname_id"`
		IsoCode    string            `maxminddb:"iso_code"`
		Names      map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
	Traits struct {
		AutonomousSystemNumber       uint    `maxminddb:"autonomous_system_number"`
		AutonomousSystemOrganization string  `maxminddb:"autonomous_system_organization"`
		ConnectionType               string  `maxminddb:"connection_type"`
		Domain                       string  `maxminddb:"domain"`
		IsAnonymousProxy             bool    `maxminddb:"is_anonymous_proxy"`
		IsLegitimateProxy            bool    `maxminddb:"is_legitimate_proxy"`
		IsSatelliteProvider          bool    `maxminddb:"is_satellite_provider"`
		ISP                          string  `maxminddb:"isp"`
		MobileCountryCode            string  `maxminddb:"mobile_country_code"`
		MobileNetworkCode            string  `maxminddb:"mobile_network_code"`
		Organization                 string  `maxminddb:"organization"`
		StaticIPScore                float64 `maxminddb:"static_ip_score"`
		UserType                     string  `maxminddb:"user_type"`
	} `maxminddb:"traits"`
}

// The City struct corresponds to the data in the GeoIP2/GeoLite2 City
// databases.
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

// The Country struct corresponds to the data in the GeoIP2/GeoLite2
// Country databases.
type Country struct {
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
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
}

// The AnonymousIP struct corresponds to the data in the GeoIP2
// Anonymous IP database.
type AnonymousIP struct {
	IsAnonymous        bool `maxminddb:"is_anonymous"`
	IsAnonymousVPN     bool `maxminddb:"is_anonymous_vpn"`
	IsHostingProvider  bool `maxminddb:"is_hosting_provider"`
	IsPublicProxy      bool `maxminddb:"is_public_proxy"`
	IsResidentialProxy bool `maxminddb:"is_residential_proxy"`
	IsTorExitNode      bool `maxminddb:"is_tor_exit_node"`
}

// The ASN struct corresponds to the data in the GeoLite2 ASN database.
type ASN struct {
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"`
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
}

// The ConnectionType struct corresponds to the data in the GeoIP2
// Connection-Type database.
type ConnectionType struct {
	ConnectionType string `maxminddb:"connection_type"`
}

// The Domain struct corresponds to the data in the GeoIP2 Domain database.
type Domain struct {
	Domain string `maxminddb:"domain"`
}

// The ISP struct corresponds to the data in the GeoIP2 ISP database.
type ISP struct {
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"`
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
	ISP                          string `maxminddb:"isp"`
	MobileCountryCode            string `maxminddb:"mobile_country_code"`
	MobileNetworkCode            string `maxminddb:"mobile_network_code"`
	Organization                 string `maxminddb:"organization"`
}

type databaseType int

const (
	isAnonymousIP = 1 << iota
	isASN
	isCity
	isConnectionType
	isCountry
	isDomain
	isEnterprise
	isISP
)

func getDBType(dbType string) databaseType {
	switch dbType {
	case "GeoIP2-Anonymous-IP":
		return isAnonymousIP
	case "DBIP-ASN-Lite (compat=GeoLite2-ASN)",
		"GeoLite2-ASN":
		return isASN
	// We allow City lookups on Country for back compat
	case "DBIP-City-Lite",
		"DBIP-Country-Lite",
		"DBIP-Country",
		"DBIP-Location (compat=City)",
		"GeoLite2-City",
		"GeoIP2-City",
		"GeoIP2-City-Africa",
		"GeoIP2-City-Asia-Pacific",
		"GeoIP2-City-Europe",
		"GeoIP2-City-North-America",
		"GeoIP2-City-South-America",
		"GeoIP2-Precision-City",
		"GeoLite2-Country",
		"GeoIP2-Country":
		return isCity | isCountry
	case "GeoIP2-Connection-Type":
		return isConnectionType
	case "GeoIP2-Domain":
		return isDomain
	case "DBIP-ISP (compat=Enterprise)",
		"DBIP-Location-ISP (compat=Enterprise)",
		"GeoIP2-Enterprise":
		return isEnterprise | isCity | isCountry
	case "GeoIP2-ISP",
		"GeoIP2-Precision-ISP":
		return isISP | isASN
	default:
		return 0
	}
}
