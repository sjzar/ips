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

package operate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjzar/ips/pkg/model"
)

func TestFieldSelector(t *testing.T) {
	ast := assert.New(t)

	meta := &model.Meta{
		Fields:     []string{"country", "province", "city", "isp"},
		FieldAlias: make(map[string]string),
	}
	selector, err := NewFieldSelector(meta, "chinaCity")
	ast.Nil(err)
	ast.Equal([]string{"country", "province", "city", "isp"}, selector.Fields())

	info := &model.IPInfo{
		Data: map[string]string{
			"country":  "中国",
			"province": "浙江",
			"city":     "杭州",
			"isp":      "电信",
		},
	}
	err = selector.Do(info)
	ast.Nil(err)
	ast.Equal([]string{"中国", "浙江", "杭州", "电信"}, info.Values())

	info.Data = map[string]string{
		"country":  "日本",
		"province": "东京都",
		"city":     "品川区",
		"isp":      "WIDE Project",
	}
	err = selector.Do(info)
	ast.Nil(err)
	ast.Equal([]string{"日本", "", "", ""}, info.Values())
}
