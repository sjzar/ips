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
	"bufio"
	"io"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sjzar/ips/internal/data"

	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

const (
	DataSep          = "\t"
	DataConditionSep = ","
)

// Data is segmented by '\t', each line is formatted as: <condition>\t<replace>\n
// @ <condition> - match condition, in `url.Values` format, supports multiple field matching
// @ <replace> - rewrite content, in `url.Values` format, supports multiple field rewrites
// Example:
// 	# match condition: province=内蒙古
// 	# rewrite content: province=内蒙
// 	province=内蒙古	province=内蒙
// 	# match condition: asn=4134
// 	# rewrite content: isp=电信
// 	asn=4134	isp=电信

// DataRewriter is responsible for rewriting data based on the provided rules.
// It uses predefined conditions to determine which data entries should be modified.
type DataRewriter struct {

	// Data maps conditions to rewrite units. Multiple fields can be matched in sequence to determine the rewrite rule.
	// level 1: map[string] list of fields, multiple fields separated by DataConditionSep, each rewrite needs to match all the fields.
	// level 2: map[string] Match condition, multiple fields separated by DataConditionSep.
	// level 3: []RewriteUnit rewrite unit, rewrites sequentially, in the order it was loaded
	Data map[string]map[string][]RewriteUnit
}

// RewriteUnit is a structure that holds the condition to match and the content to replace.
type RewriteUnit struct {
	Condition url.Values // Condition to match for rewriting
	Replace   url.Values // Content to replace if the condition is matched
}

// NewDataRewriter initializes and returns a new DataRewriter.
func NewDataRewriter() *DataRewriter {
	return &DataRewriter{
		Data: make(map[string]map[string][]RewriteUnit),
	}
}

// Do applies the data rewriting rules to the provided IPInfo.
func (d *DataRewriter) Do(info *model.IPInfo) error {
	if d.Data == nil || len(d.Data) == 0 {
		return nil
	}
	for fieldsStr, match := range d.Data {
		fields := strings.Split(fieldsStr, DataConditionSep)
		values := make([]string, 0, len(fields))
		for _, field := range fields {
			value, _ := info.GetData(field)
			values = append(values, value)
		}
		valuesStr := strings.Join(values, DataConditionSep)
		replace, ok := match[valuesStr]
		if !ok {
			continue
		}
		for _, unit := range replace {
			for field, value := range unit.Replace {
				if _, ok := info.FieldAlias[field]; ok {
					field = info.FieldAlias[field]
				}
				info.Data[field] = value[0]
			}
		}
	}
	return nil
}

// LoadFile loads the rewriting rules from a file.
func (d *DataRewriter) LoadFile(file string) error {
	if len(file) == 0 {
		return errors.ErrFileEmpty
	}

	if _, ok := data.PresetFiles[file]; ok {
		d.LoadString(data.PresetFiles[file])
		return nil
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	return d.load(f)
}

// LoadFiles loads the rewriting rules from a list of files.
func (d *DataRewriter) LoadFiles(files []string) error {
	for _, file := range files {
		if err := d.LoadFile(file); err != nil {
			return err
		}
	}
	return nil
}

// LoadString loads rewriting rules from the provided data strings.
func (d *DataRewriter) LoadString(data ...string) {
	_ = d.load(strings.NewReader(strings.Join(data, "\n")))
}

// load is a utility function that reads rewriting rules from an io.Reader.
// It parses the conditions and the corresponding rewriting rules.
func (d *DataRewriter) load(r io.Reader) error {
	if d.Data == nil {
		d.Data = make(map[string]map[string][]RewriteUnit)
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		split := strings.SplitN(string(line), DataSep, 2)
		if len(split) < 2 {
			continue
		}
		condition, err := url.ParseQuery(split[0])
		if err != nil {
			log.Errorf("[%s] parse condition failed: %v", split[0], err)
			return err
		}
		fields := make([]string, 0, len(condition))
		values := make([]string, 0, len(condition))
		for k, v := range condition {
			fields = append(fields, k)
			values = append(values, v[0])
		}
		fieldsStr := strings.Join(fields, DataConditionSep)
		valuesStr := strings.Join(values, DataConditionSep)

		replace, err := url.ParseQuery(split[1])
		if err != nil {
			log.Errorf("[%s] parse replace failed: %v", split[1], err)
			return err
		}

		if d.Data[fieldsStr] == nil {
			d.Data[fieldsStr] = make(map[string][]RewriteUnit)
		}
		if d.Data[fieldsStr][valuesStr] == nil {
			d.Data[fieldsStr][valuesStr] = make([]RewriteUnit, 0)
		}

		d.Data[fieldsStr][valuesStr] = append(d.Data[fieldsStr][valuesStr], RewriteUnit{
			Condition: condition,
			Replace:   replace,
		})
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
