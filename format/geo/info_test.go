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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo_Name(t *testing.T) {
	ast := assert.New(t)

	info := &Info{
		Names: map[string]string{
			"en":    "hello",
			"zh-CN": "你好",
		},
	}

	// Test English name
	ast.Equal("hello", info.Name("en"), "Expected English name not found")

	// Test Chinese name
	ast.Equal("你好", info.Name("zh-CN"), "Expected Chinese name not found")

	// Test fallback to default language (zh-CN)
	ast.Equal("你好", info.Name("es"), "Expected fallback to default language name not found")

	// Test fallback to English name, default language (es) is not found
	ast.Nil(SetLanguage("es"))
	ast.Equal("hello", info.Name("es"), "Expected fallback to English name not found")

}
func TestParseGeoInfo(t *testing.T) {
	tests := []struct {
		input  string
		output *Info
		ok     bool
	}{
		{
			"123\ten:hello|zh-CN:你好\tcode\tiso_code\t1\t1\t2\t3|4",
			&Info{
				GeoNameID: 123,
				Names: map[string]string{
					"en":    "hello",
					"zh-CN": "你好",
				},
				Code:              "code",
				IsoCode:           "iso_code",
				IsInEuropeanUnion: true,
				ConnectionID:      1,
				CountryID:         2,
				SubdivisionIDs:    []int{3, 4},
			},
			true,
		},
		{
			"invalid_data",
			nil,
			false,
		},
	}

	for _, tt := range tests {
		result, ok := ParseGeoInfo(tt.input)
		assert.Equal(t, tt.ok, ok, "Expected parse status not found")
		if ok {
			assert.Equal(t, tt.output, result, "Expected parsed Info not found")
		}
	}
}
