package reader

import (
	"log"
	"strings"

	"github.com/sjzar/ips/db"
	"github.com/sjzar/ips/model"
)

// magicArgs 魔法参数
// 用于快速匹配过滤字段
var magicArgs = map[string]string{
	"":               "country,province,city,isp",
	"fullFields":     strings.Join(model.FullFields, ","),
	"awdbFullFields": strings.Join(db.AWDBFullFields, ","),
	"ipdbFullFields": strings.Join(db.IPDBFullFields, ","),
	"chinaCity":      "country#中国:country,province,city,isp|country,-,-,-",
	"provinceAndISP": "country#中国:province,isp|-,-",
}

// KeyFieldsFilter 字段过滤工具
type KeyFieldsFilter struct {
	key           string
	defaultFields []string
	keyFields     map[string][]string
}

// NewKeyFieldsFilter 初始化字段过滤器
// 语法: <key>#<value1>:<filed1>,<field2>,<field3>|<value2>:<filed1>,<field2>,<field3>
//     - <key>: 用于匹配value的字段，默认为country
//     - <value>: 用于匹配fields的值，允许缺省，缺省时表示默认fields
//     - <filed>: 输出的字段内容
// 举例: country#中国:country,province,city,isp|country,-,-,isp
//     - 当country为中国时，输出country,province,city,isp；country为其他时，输出 country,-,-,isp (-表示缺省，不输出)
// 约束: 每组的字段数量需要保持一致
func NewKeyFieldsFilter(arg string) *KeyFieldsFilter {
	if _arg, ok := magicArgs[arg]; ok {
		return NewKeyFieldsFilter(_arg)
	}

	filter := &KeyFieldsFilter{
		key:           model.Country,
		defaultFields: make([]string, 0),
		keyFields:     make(map[string][]string),
	}

	split := strings.SplitN(arg, "#", 2)
	if len(split) == 2 {
		filter.key = split[0]
		split[0] = split[1]
	}

	groups := strings.Split(split[0], "|")
	fieldsNum := 0
	for _, group := range groups {
		keyFields := strings.SplitN(group, ":", 2)
		key := ""
		if len(keyFields) == 2 {
			key = keyFields[0]
			keyFields[0] = keyFields[1]
		}
		fields := strings.Split(keyFields[0], ",")
		if fieldsNum == 0 {
			fieldsNum = len(fields)
		}
		if len(fields) == 0 || len(fields) != fieldsNum {
			log.Fatalf("invalid key fields: [%s]\n", arg)
		}

		filter.keyFields[key] = fields
	}

	if len(filter.keyFields) == 0 {
		return NewKeyFieldsFilter("")
	}

	if defaultFields, ok := filter.keyFields[""]; ok {
		filter.defaultFields = defaultFields
	} else {
		for key := range filter.keyFields {
			filter.defaultFields = filter.keyFields[key]
			break
		}
	}

	return filter
}

// Fields 查询字段列表
func (f *KeyFieldsFilter) Fields(key string) []string {
	if fields, ok := f.keyFields[key]; ok {
		return fields
	}
	return f.defaultFields
}

// Key 条件字段
func (f *KeyFieldsFilter) Key() string {
	return f.key
}

// FieldsCollection 字段集合
func (f *KeyFieldsFilter) FieldsCollection() string {
	collection := make([]string, len(f.defaultFields))
	placeholderIndex := make([]int, 0)
	for i := range f.defaultFields {
		if f.defaultFields[i] == model.Placeholder {
			placeholderIndex = append(placeholderIndex, i)
			continue
		}
		collection[i] = f.defaultFields[i]
	}
	if len(placeholderIndex) != 0 {
		for _, i := range placeholderIndex {
			for _, fields := range f.keyFields {
				if fields[i] != model.Placeholder {
					collection[i] = fields[i]
					break
				}
			}
		}
	}

	return strings.Join(collection, ",")
}
