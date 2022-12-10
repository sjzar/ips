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

package rewriter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var DataLoad1 = `isp	Beihai	电信
province	上海市	上海
province	云南省	云南
[wrong line]
asnumber	4808	isp|联通
`

func TestDataRewriter(t *testing.T) {
	ast := assert.New(t)

	dataLoader := NewDataLoader()
	dataLoader.LoadString(DataLoad1)
	dataRewriter := NewDataRewriter(dataLoader, nil)
	data := map[string]string{}
	_, _, retData, err := dataRewriter.Rewrite(nil, nil, data)
	ast.Nil(err)
	ast.Equal(0, len(retData))

	data = map[string]string{
		"province": "上海市",
		"isp":      "Beihai",
	}
	_, _, retData, err = dataRewriter.Rewrite(nil, nil, data)
	ast.Nil(err)
	ast.Equal(2, len(retData))
	ast.Equal("上海", retData["province"])
	ast.Equal("电信", retData["isp"])

	data = map[string]string{
		"asnumber": "4808",
	}

	_, _, retData, err = dataRewriter.Rewrite(nil, nil, data)
	ast.Nil(err)
	ast.Equal(2, len(retData))
	ast.Equal("4808", retData["asnumber"])
	ast.Equal("联通", retData["isp"])
}
