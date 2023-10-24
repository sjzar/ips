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

func TestTranslate(t *testing.T) {
	ast := assert.New(t)

	err := SetLanguage(LangEnglish)
	ast.Nil(err)

	ast.Equal("Asia", Translate("continent", "亚洲"))
	ast.Equal("China", Translate("country", "中国"))
	ast.Equal("China", Translate("", "中国"))
	ast.Equal("Shanghai", Translate("province", "上海"))
	ast.Equal("Shanghai", Translate("region", "上海"))
	ast.Equal("Shanghai", Translate("city", "上海"))
	ast.Equal("Los Angeles", Translate("city", "洛杉矶"))
	ast.Equal("England", Translate("province", "英格兰"))
}
