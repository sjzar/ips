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
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

// The format is: <fields>[|<rule1>:<fields>|<rule2>:<fields>|<default fields>]
// @ <fields> - a list of fields to output, separated by ","
// @ <rule> - match rules, expressed as url.Values
// Example:
//  # select fields: country,province,city,isp
//  # if country is 中国, select fields: country, province, city, isp
//  # if country is not 中国, select fields: country, <empty value>, <empty value>, <empty value>
// 	"country,province,city,isp|country=!中国:country"

// Constants used as separators and rules in field selectors.
const (
	SelectorGroupSep    = "|"
	SelectorFieldSep    = ","
	SelectorRuleSep     = ":"
	SelectorValueNot    = "!"
	SelectorValueOr     = "/"
	SelectorValReplace  = "="
	SelectorWildcardArg = "*"
)

// magicArgs is a map that defines shortcuts for commonly used selector arguments.
var magicArgs = map[string]string{
	"find":           "country,province,city,isp",
	"chinaCity":      "country,province,city,isp|country=!中国:country",
	"provinceAndISP": "province,isp|country=!中国:",
	"cn":             "country|country=中国:country='CN'|country='OV'",
}

// NewFieldSelector initializes a FieldSelector based on the given argument string.
// It supports a custom string format to determine which fields to select and under which conditions.
// The format is: <fields>[|<rule1>:<fields>|<rule2>:<fields>|<default fields>]
// MAKE SURE that the incoming meta comes from StandardReader, not DB Reader, since it may modify the *model.Meta.
// Refer to the provided examples in the function comments for more details.
func NewFieldSelector(meta *model.Meta, arg string) (*FieldSelector, error) {

	// Check if the arg is a magic argument and retrieve its actual value.
	if _arg, ok := magicArgs[arg]; ok {
		return NewFieldSelector(meta, _arg)
	}

	// Ensure that meta has fields defined.
	if meta.Fields == nil {
		return nil, errors.ErrMetaFieldsUndefined
	}

	ret := &FieldSelector{
		fields: meta.Fields,
		rules:  make([]*FieldSelectorRule, 0),
	}

	// If the arg is empty or a wildcard, return the selector as is.
	if len(arg) == 0 || arg == SelectorWildcardArg {
		return ret, nil
	}

	var err error
	metaFields := meta.SupportFields()
	groups := strings.Split(arg, SelectorGroupSep)
	if ret.fields, err = processFields(groups[0], metaFields, meta.FieldAlias); err != nil {
		return nil, err
	}

	// Update the meta fields.
	meta.Fields = ret.fields

	if len(groups) == 1 {
		return ret, nil
	}

	for _, groupStr := range groups[1:] {
		rule, err := processRule(groupStr, ret.fields, meta.FieldAlias)
		if err != nil {
			return nil, err
		}
		ret.rules = append(ret.rules, rule)
	}

	return ret, nil
}

// Do applies the field selector rules to an IPInfo object, selecting and modifying fields as necessary.
func (f *FieldSelector) Do(info *model.IPInfo) error {
	fields := f.fields
	replaceFields := make(map[string]string)
	for _, rule := range f.rules {
		if rule.IsMatch(info) {
			fields = rule.Fields
			replaceFields = rule.ReplaceFields
			break
		}
	}

	info.Fields = fields
	info.ReplaceFields = replaceFields
	return nil
}

// processFields validates and processes a field string based on available meta fields and field aliases.
func processFields(fieldStr string, metaFields map[string]bool, fieldAlias map[string]string) ([]string, error) {
	fields := strings.Split(fieldStr, SelectorFieldSep)
	for index, field := range fields {
		if _, ok := fieldAlias[field]; ok {
			fields[index] = fieldAlias[field]
		}
	}
	return fields, nil
}

// processRule processes a rule string and constructs a FieldSelectorRule from it.
func processRule(ruleStr string, fullFields []string, fieldAlias map[string]string) (*FieldSelectorRule, error) {
	var err error
	parts := strings.SplitN(ruleStr, SelectorRuleSep, 2)

	condition := url.Values{}
	fieldStr := parts[0]
	if len(parts) == 2 {
		condition, err = url.ParseQuery(parts[0])
		if err != nil {
			return nil, err
		}
		fieldStr = parts[1]
	}

	fields, replaceFields, err := fillFields(fullFields, strings.Split(fieldStr, SelectorFieldSep), fieldAlias)
	if err != nil {
		return nil, err
	}

	return &FieldSelectorRule{
		Condition:     condition,
		Fields:        fields,
		ReplaceFields: replaceFields,
	}, nil
}

// fillFields fills in missing fields based on fullFields and processes field replacements.
func fillFields(fullFields []string, fields []string, fieldAlias map[string]string) ([]string, map[string]string, error) {

	// Create a set for fast lookup
	tmpFullFields := make(map[string]bool)
	for _, field := range fullFields {
		tmpFullFields[field] = true
	}

	// Process fields and create a set for fields to be filled
	tmpFields := make(map[string]bool)
	replaceFields := make(map[string]string)
	for _, field := range fields {
		if strings.Contains(field, SelectorValReplace) {
			split := strings.SplitN(field, SelectorValReplace, 2)
			replaceFields[split[0]] = split[1]
			field = split[0]
		}
		if _, ok := tmpFullFields[field]; !ok {
			log.Debugf("invalid field: %s", field)
			return nil, nil, errors.ErrFieldInvalid
		}
		tmpFields[field] = true
	}

	ret := make([]string, len(fullFields))
	for i, field := range fullFields {
		if _, ok := tmpFields[field]; ok {
			alias, hasAlias := fieldAlias[field]
			if hasAlias {
				ret[i] = alias
			} else {
				ret[i] = field
			}
		}
	}

	return ret, replaceFields, nil
}

// Fields returns the list of fields currently set in the field selector.
func (f *FieldSelector) Fields() []string {
	return f.fields
}

// FieldSelector represents a selector that can be applied to IPInfo to select and modify fields.
type FieldSelector struct {
	fields []string
	rules  []*FieldSelectorRule
}

// FieldSelectorRule represents a single rule within a FieldSelector.
type FieldSelectorRule struct {
	Condition     url.Values
	Fields        []string
	ReplaceFields map[string]string
}

// IsMatch determines if the given IPInfo matches the conditions set in the rule.
func (r *FieldSelectorRule) IsMatch(info *model.IPInfo) bool {
	for key, values := range r.Condition {
		actual, _ := info.GetData(key)
		for _, value := range values {
			not := false
			if strings.HasPrefix(value, SelectorValueNot) {
				not = true
				value = strings.TrimPrefix(value, SelectorValueNot)
			}

			match := false
			for _, v := range strings.Split(value, SelectorValueOr) {
				if v == actual {
					match = true
					break
				}
			}

			// match || not && !match
			if match == not {
				return false
			}
		}
	}
	return true
}
