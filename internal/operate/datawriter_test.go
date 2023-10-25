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

var DataLoad1 = `isp=Beihai	isp=电信
province=上海市	province=上海
province=云南省	province=云南
[wrong line]
asnumber=4808	isp=联通&country=中国
`

func TestDataRewriter(t *testing.T) {
	ast := assert.New(t)

	dataRewriter := NewDataRewriter()
	dataRewriter.LoadString(DataLoad1)

	info := &model.IPInfo{
		Data: map[string]string{},
	}

	err := dataRewriter.Do(info)
	ast.Nil(err)
	ast.Equal(0, len(info.Data))

	info = &model.IPInfo{
		Data: map[string]string{
			"province": "上海市",
			"isp":      "Beihai",
		},
	}
	err = dataRewriter.Do(info)
	ast.Nil(err)
	ast.Equal(2, len(info.Data))
	value, _ := info.GetData("province")
	ast.Equal("上海", value)
	value, _ = info.GetData("isp")
	ast.Equal("电信", value)

	info = &model.IPInfo{
		Data: map[string]string{
			"asnumber": "4808",
		},
	}

	err = dataRewriter.Do(info)
	ast.Nil(err)
	ast.Equal(3, len(info.Data))
	value, _ = info.GetData("asnumber")
	ast.Equal("4808", value)
	value, _ = info.GetData("isp")
	ast.Equal("联通", value)
	value, _ = info.GetData("country")
	ast.Equal("中国", value)
}
