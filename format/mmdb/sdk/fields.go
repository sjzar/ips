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

package sdk

// PreProcessFieldAlias maps fields from various sources to a standardized format.
// This is useful for normalizing field names from different IP information providers.
var PreProcessFieldAlias = map[string]string{
	// GeoIP2 like
	"location_accuracy_radius":              "accuracy_radius",
	"location_latitude":                     "latitude",
	"location_longitude":                    "longitude",
	"location_metro_code":                   "metro_code",
	"location_time_zone":                    "time_zone",
	"traits_autonomous_system_number":       "autonomous_system_number",
	"traits_autonomous_system_organization": "autonomous_system_organization",
	"traits_connection_type":                "connection_type",
	"traits_domain":                         "domain",
	"traits_is_anonymous_proxy":             "is_anonymous_proxy",
	"traits_is_legitimate_proxy":            "is_legitimate_proxy",
	"traits_is_satellite_provider":          "is_satellite_provider",
	"traits_isp":                            "isp",
	"traits_mobile_country_code":            "mobile_country_code",
	"traits_mobile_network_code":            "mobile_network_code",
	"traits_organization":                   "organization",
	"traits_static_ip_score":                "static_ip_score",
	"traits_user_type":                      "user_type",

	// IPInfo.io
	"continent_name": "continent",
	"country_name":   "country",
	"asn":            "autonomous_system_number",
	"as_name":        "autonomous_system_organization",
	"as_domain":      "domain",
}

// DoPreProcessFieldAlias takes a map of data fields and standardizes their names
// using the PreProcessFieldAlias mapping. This helps ensure consistency in field names
// regardless of the data source.
func DoPreProcessFieldAlias(data map[string]string) map[string]string {
	for k, v := range PreProcessFieldAlias {
		if _, ok := data[k]; ok {
			data[v] = data[k]
			delete(data, k)
		}
	}
	return data
}
