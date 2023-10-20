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

package model

import (
	"net"
	"strings"

	"github.com/sjzar/ips/ipnet"
)

// IPInfo represents information related to a specific IP address.
type IPInfo struct {

	// IP represents the IP address.
	IP net.IP

	// IPNet denotes the network range to which the IP belongs.
	IPNet *ipnet.Range

	// Data holds the actual information related to the IP in a key-value format.
	Data map[string]string

	// FieldAlias maps common field names to their corresponding database field names.
	FieldAlias map[string]string

	// Fields lists the fields that should be output.
	Fields []string

	// ReplaceFields specifies fields that should be replaced with specific values.
	// operate.FieldSelector use this to replace the value of the field.
	ReplaceFields map[string]string
}

// GetData retrieves the data for the given field.
// It first checks if the field exists in the Data map directly. If not, it checks
// the FieldAlias map to determine the corresponding database field and retrieves the value.
// Returns the value and a boolean indicating if the field was found.
func (i *IPInfo) GetData(field string) (string, bool) {
	if v, ok := i.Data[field]; ok {
		return v, true
	}
	if dbField, ok := i.FieldAlias[field]; ok {
		if v, ok := i.Data[dbField]; ok {
			return v, true
		}
	}
	return "", false
}

// AddCommonFieldAlias adds common field aliases to the FieldAlias map.
// It maps a database-specific field name to a common field name.
func (i *IPInfo) AddCommonFieldAlias(fieldAlias map[string]string) {
	if i.FieldAlias == nil {
		i.FieldAlias = make(map[string]string)
	}

	for commonField, dbField := range fieldAlias {
		if _, ok := i.Data[dbField]; ok {
			i.FieldAlias[commonField] = dbField
		}
	}
}

// Values returns a slice of data values based on the Fields slice.
// It also considers any replacements specified in the ReplaceFields map.
func (i *IPInfo) Values() []string {
	ret := make([]string, len(i.Fields))
	for index, field := range i.Fields {
		if v, ok := i.GetData(field); ok {
			ret[index] = v
		}
		if replaceVal, ok := i.ReplaceFields[field]; ok {
			ret[index] = strings.Trim(replaceVal, "'")
		}
	}
	return ret
}

// IPInfoOutput represents the structure for outputting IP information.
type IPInfoOutput struct {
	IP   string            `json:"ip"`
	Net  string            `json:"net"`
	Data map[string]string `json:"data"`
}

// Output constructs and returns an IPInfoOutput based on the current IPInfo.
// It decides whether to use the database field or the common field based on the dbFiled flag.
func (i *IPInfo) Output(dbFiled bool) *IPInfoOutput {
	data := make(map[string]string, len(i.Fields))
	values := i.Values()

	fieldAliasReverse := make(map[string]string, len(i.FieldAlias))
	for commonField, dbField := range i.FieldAlias {
		fieldAliasReverse[dbField] = commonField
	}
	for index, field := range i.Fields {
		if dbFiled {
			data[field] = values[index]
			continue
		}

		if commonField, ok := fieldAliasReverse[field]; ok {
			data[commonField] = values[index]
			continue
		}
		data[field] = values[index]
	}

	ipNet := ""
	ipNets := i.IPNet.IPNets()
	if len(ipNets) > 0 {
		ipNet = ipNets[0].String()
	}

	return &IPInfoOutput{
		IP:   i.IP.String(),
		Net:  ipNet,
		Data: data,
	}
}
