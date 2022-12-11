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
	"strings"

	"github.com/sjzar/ips/db/awdb"
	"github.com/sjzar/ips/db/ipdb"
	"github.com/sjzar/ips/model"
)

const (
	SelectorGroupSep    = "|"
	SelectorFieldSep    = ","
	SelectorRuleSep     = ":"
	SelectorKeyValueSep = "="
)

// magicArgs 魔法参数
// 用于快速匹配选择器
var magicArgs = map[string]string{
	"":               "country,province,city,isp",
	"full":           strings.Join(model.CommonFields, ","),
	"awdb":           strings.Join(awdb.FullFields, ","),
	"ipdb":           strings.Join(ipdb.FullFields, ","),
	"chinaCity":      "country,province,city,isp|country=中国:country,province,city,isp|country,,,",
	"provinceAndISP": "province,isp|country=中国:province,isp|,",
	"cn":             "country|country=中国:'CN'|'OV'",
}

// FieldSelector 字段选择器
type FieldSelector struct {
	fields        []string
	defaultFields []string
	rules         map[string]map[string][]string
}

// NewFieldSelector 初始化字段选择器
// 用于选择输出的字段，支持简单的匹配逻辑
// 输入的参数格式为: <fields>[|<rule1>:<fields>|<rule2>:<fields>|<default fields>]
// @ <fields> - 字段列表，表示需要输出的字段，字段之间使用","分隔
// @ <rule> - 匹配规则，非必填项，以<key>=<value>表示，匹配上的话，则使用<rule>对应的<fields>，匹配优先级为<fields>的顺序
// @ <default fields> - 默认字段列表，非必填项，若<default fileds>为空，则使用<fields>
// @ sep - Select函数输出数据的分隔符，默认为","
// 举例
// country,province,city,isp|country=中国:country,province,city,isp|country,-,-,-
// 针对国家区分IP库精度，若国家是"中国"，返回"国家,省份,城市,运营商"，若国家不是"中国"，返回"国家,-,-,-"
// country|country=中国:CN|OV
// 返回中国和海外
// isp|isp=电信:电信|isp=联通:联通|isp=移动:移动|其他
// 返回电信、联通、移动和其他运营商
func NewFieldSelector(arg string) *FieldSelector {
	if _arg, ok := magicArgs[arg]; ok {
		return NewFieldSelector(_arg)
	}
	ret := &FieldSelector{
		rules: make(map[string]map[string][]string),
	}
	groups := strings.Split(arg, SelectorGroupSep)
	for i, group := range groups {
		fields := strings.Split(group, SelectorFieldSep)
		if i == 0 {
			ret.fields = fields
			continue
		}
		if len(fields) != len(ret.fields) {
			continue
		}
		rule := strings.Split(fields[0], SelectorRuleSep)
		if len(rule) == 1 {
			ret.defaultFields = fields
			continue
		}
		keyValue := strings.Split(rule[0], SelectorKeyValueSep)
		if len(keyValue) == 2 {
			fields[0] = rule[1]
			valueFields, ok := ret.rules[keyValue[0]]
			if !ok {
				valueFields = make(map[string][]string)
				ret.rules[keyValue[0]] = valueFields
			}
			valueFields[keyValue[1]] = fields
		}
	}
	return ret
}

// Select 选择字段
func (f *FieldSelector) Select(data map[string]string) []string {
	fields := f.fields
	if len(f.defaultFields) != 0 {
		fields = f.defaultFields
	}

	for _, field := range f.fields {
		if valueFields, ok := f.rules[field]; ok && len(valueFields[data[field]]) != 0 {
			fields = valueFields[data[field]]
		}
	}

	ret := make([]string, len(fields))
	for i, field := range fields {
		trimField := strings.TrimSuffix(strings.TrimPrefix(field, "'"), "'")
		if val, ok := data[trimField]; ok {
			ret[i] = val
		} else if len(field) > len(trimField) {
			ret[i] = trimField
		}
	}
	return ret
}

// Fields 返回字段列表
func (f *FieldSelector) Fields() []string {
	return f.fields
}
