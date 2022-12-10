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

package ipio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldSelector(t *testing.T) {
	ast := assert.New(t)

	selector := NewFieldSelector("chinaCity")
	ast.Equal([]string{"country", "province", "city", "isp"}, selector.Fields())

	data := map[string]string{
		"country":  "中国",
		"province": "浙江",
		"city":     "杭州",
		"isp":      "电信",
	}
	ast.Equal([]string{"中国", "浙江", "杭州", "电信"}, selector.Select(data))

	data = map[string]string{
		"country":  "日本",
		"province": "东京都",
		"city":     "品川区",
		"isp":      "WIDE Project",
	}
	ast.Equal([]string{"日本", "", "", ""}, selector.Select(data))
}
