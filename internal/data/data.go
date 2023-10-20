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

package data

import (
	_ "embed"
)

// PresetFiles is a map of preset data files.
var PresetFiles = map[string]string{
	"asn2isp":       ASN2ISP,
	"city":          City,
	"isp":           ISP,
	"province":      Province,
	"qqwry_area":    QQwryArea,
	"qqwry_country": QQwryCountry,
}

//go:embed asn2isp.map
var ASN2ISP string

//go:embed city.map
var City string

//go:embed isp.map
var ISP string

//go:embed province.map
var Province string

//go:embed qqwry_area.map
var QQwryArea string

//go:embed qqwry_country.map
var QQwryCountry string
