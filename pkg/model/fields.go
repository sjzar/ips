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

	// ASN 自治域号
	ASN = "asn"

	// Continent 大洲
	Continent = "continent"

	// UTCOffset UTC 偏移值
	UTCOffset = "utcOffset"

	// Latitude 纬度
	Latitude = "latitude"

	// Longitude 经度
	Longitude = "longitude"

	// ChinaAdminCode 中国行政区划代码
	ChinaAdminCode = "chinaAdminCode"
)

// ConvertToDBFields converts field names from a reader database's alias to a writer database's alias.
func ConvertToDBFields(fields []string, readerFieldAlias, writerFieldAlias map[string]string) []string {

	// Create a reverse mapping for the reader's alias for quick lookup.
	readerFieldAliasReverse := make(map[string]string)
	for commonField, dbField := range readerFieldAlias {
		readerFieldAliasReverse[dbField] = commonField
	}

	// Iterate through the fields and perform the conversion.
	convertedFields := make([]string, len(fields))
	for i, field := range fields {
		// If the field exists in the reader's reverse alias, use the standard field name.
		if standardField, exists := readerFieldAliasReverse[field]; exists {
			field = standardField
		}
		// If the field exists in the writer's alias, use the writer's alias.
		if writerField, exists := writerFieldAlias[field]; exists {
			field = writerField
		}
		convertedFields[i] = field
	}

	return convertedFields
}
