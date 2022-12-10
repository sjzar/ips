/*
 * Copyright (c) 2022 shenjunzheng@gmail.com
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

	// Continent 大洲
	Continent = "continent"

	// UTCOffset UTC偏移值
	UTCOffset = "utcOffset"

	// Latitude 纬度
	Latitude = "latitude"

	// Longitude 经度
	Longitude = "longitude"

	// ChinaAdminCode 中国行政区划代码
	ChinaAdminCode = "chinaAdminCode"
)

// CommonFields 公共字段
var CommonFields = []string{
	Country,
	Province,
	City,
	ISP,
	Continent,
	UTCOffset,
	Latitude,
	Longitude,
	ChinaAdminCode,
}
