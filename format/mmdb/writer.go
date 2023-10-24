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

package mmdb

import (
	"io"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"

	"github.com/sjzar/ips/format/geo"
	"github.com/sjzar/ips/pkg/errors"
	"github.com/sjzar/ips/pkg/model"
)

// Writer provides functionalities to write IP data into MMDB format.
type Writer struct {
	meta   *model.Meta      // Metadata for the IP database
	writer *mmdbwriter.Tree // MMDB writer instance
	option WriterOption     // Writer options
}

// NewWriter initializes a new Writer instance for writing IP data in MMDB format.
func NewWriter(meta *model.Meta) (*Writer, error) {
	opts := mmdbwriter.Options{
		DatabaseType: "GeoIP2-City",
	}

	writer, err := mmdbwriter.New(opts)
	if err != nil {
		return nil, err
	}

	return &Writer{
		meta:   meta,
		writer: writer,
	}, nil
}

// WriterOption provides options for the Writer.
type WriterOption struct {
	SelectLanguages string // SelectLanguages specifies the languages to be selected for the names.
}

// SetOption sets the provided options to the Writer.
// Currently, it supports mmdbwriter.Options for the MMDB writer.
func (w *Writer) SetOption(option interface{}) error {
	if opt, ok := option.(WriterOption); ok {
		w.option = opt
		return nil
	}
	if opts, ok := option.(mmdbwriter.Options); ok {
		writer, err := mmdbwriter.New(opts)
		if err != nil {
			return err
		}
		w.writer = writer
	}
	return nil
}

// Insert adds the given IP information into the writer.
func (w *Writer) Insert(info *model.IPInfo) error {
	fields := model.ConvertToDBFields(w.meta.Fields, w.meta.FieldAlias, CommonFieldsAlias)
	values := info.Values()
	if len(values) != len(fields) {
		return errors.ErrMismatchedFieldsLength
	}

	data, err := ConvertToMMDBType(w.ConvertMap(fields, values))
	if err != nil {
		return err
	}

	for _, ipNet := range info.IPNet.IPNets() {
		if err := w.insertIPNet(ipNet, data); err != nil {
			return err
		}
	}
	return nil
}

// insertIPNet inserts a single IP network into the writer.
func (w *Writer) insertIPNet(ipNet *net.IPNet, data mmdbtype.DataType) error {
	_, network, err := net.ParseCIDR(ipNet.String())
	if err != nil || network == nil {
		return nil
	}

	err = w.writer.Insert(network, data)
	if err != nil && (strings.Contains(err.Error(), "which is in a reserved network") ||
		strings.Contains(err.Error(), "which is in an aliased network")) {
		return nil
	}

	return err
}

// WriteTo writes the IP data into the provided writer in MMDB format.
func (w *Writer) WriteTo(iw io.Writer) (int64, error) {
	return w.writer.WriteTo(iw)
}

// ConvertMap converts fields and values to a map.
func (w *Writer) ConvertMap(fields, values []string) map[string]interface{} {
	ret := make(map[string]interface{})
	for i := range fields {
		value := values[i]
		if len(value) == 0 {
			continue
		}
		var convertedValue interface{}
		switch fields[i] {
		case FieldCity, FieldContinent, FieldCountry, FieldRegisteredCountry, FieldRepresentedCountry:
			convertedValue = w.convertGeoInfo(fields[i], value)
		case FieldSubdivisions:
			convertedValue = w.convertSubdivisions(value)
		case FieldAccuracyRadius, FieldMetroCode, FieldLatitude, FieldLongitude, FieldTimeZone:
			ret = w.convertLocation(fields[i], value, ret)
		case FieldPostalCode:
			convertedValue = map[string]interface{}{
				"code": value,
			}
		case FieldIsAnonymousProxy, FieldIsSatelliteProvider, "is_legitimate_proxy", FieldAutonomousSystemNumber:
			ret = w.convertTraits(fields[i], value, ret)
		default:
			convertedValue = value
		}
		if convertedValue != nil {
			ret[fields[i]] = convertedValue
		}
	}
	return ret
}

// convertGeoInfo converts the given value to its corresponding geo information.
func (w *Writer) convertGeoInfo(field, value string) interface{} {
	info, ok := geo.GetInfoByName(field, value)
	if !ok {
		return nil
	}
	return info.Map(w.option.SelectLanguages)
}

// convertLocation handles the conversion for location related fields.
func (w *Writer) convertLocation(field, value string, ret map[string]interface{}) map[string]interface{} {
	dataMap, _ := getOrCreateMap(ret, "location")
	switch field {
	case FieldAccuracyRadius, FieldMetroCode:
		if parsedValue, err := strconv.ParseUint(value, 10, 16); err == nil {
			dataMap[field] = parsedValue
		}
	case FieldLatitude, FieldLongitude:
		if parsedValue, err := strconv.ParseFloat(value, 64); err == nil {
			dataMap[field] = parsedValue
		}
	case FieldTimeZone:
		dataMap[field] = value
	}
	return ret
}

// convertTraits handles the conversion for trait related fields.
func (w *Writer) convertTraits(field, value string, ret map[string]interface{}) map[string]interface{} {
	dataMap, _ := getOrCreateMap(ret, "traits")
	switch field {
	case FieldIsAnonymousProxy, FieldIsSatelliteProvider, "is_legitimate_proxy":
		if parsedValue, err := strconv.ParseBool(value); err == nil {
			dataMap[field] = parsedValue
		}
	case FieldAutonomousSystemNumber:
		if parsedValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			dataMap[field] = mmdbtype.Uint64(parsedValue)
		}
	}
	return ret
}

// getOrCreateMap retrieves or creates a map for a given key in the provided map.
func getOrCreateMap(data map[string]interface{}, key string) (map[string]interface{}, bool) {
	if v, ok := data[key]; ok {
		if _, ok := v.(map[string]interface{}); !ok {
			return nil, false
		}
		return v.(map[string]interface{}), true
	}
	newMap := make(map[string]interface{})
	data[key] = newMap
	return newMap, true
}

// convertSubdivisions converts the given value to its corresponding subdivisions information.
func (w *Writer) convertSubdivisions(value string) []map[string]interface{} {
	split := strings.Split(value, ",")
	data := make([]map[string]interface{}, 0, len(split))
	for _, v := range split {
		info := w.convertGeoInfo("subdivisions", v)
		if info != nil {
			data = append(data, info.(map[string]interface{}))
		}
	}
	return data
}

// ConvertToMMDBType converts a given value to an appropriate MMDB data type.
// If the value is a complex type (e.g., struct or map), it uses reflection to handle the conversion.
func ConvertToMMDBType(value interface{}) (mmdbtype.DataType, error) {
	if value == nil {
		return nil, nil
	}

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return convertToMMDBType(val)
}

// convertToMMDBType converts a given value to an appropriate MMDB data type.
// If the value is a complex type (e.g., struct or map), it uses reflection to handle the conversion.
func convertToMMDBType(val reflect.Value) (mmdbtype.DataType, error) {
	switch val.Kind() {
	case reflect.Struct:
		return handleStruct(val)
	case reflect.Ptr, reflect.Interface:
		return convertToMMDBType(val.Elem())
	case reflect.Map:
		return handleMap(val)
	case reflect.Slice, reflect.Array:
		return handleSliceAndArray(val)
	default:
		return handleSimpleType(val), nil
	}
}

// handleStruct processes the provided reflect.Value of a struct type,
// and converts it into a MMDB map type.
func handleStruct(val reflect.Value) (mmdbtype.DataType, error) {
	data := make(mmdbtype.Map)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		parseValue, err := convertToMMDBType(field)
		if err != nil {
			return nil, err
		}
		data[mmdbtype.String(fieldType.Name)] = parseValue
	}
	return data, nil
}

// handleMap processes the provided reflect.Value of a map type,
// and converts it into a MMDB map type.
func handleMap(val reflect.Value) (mmdbtype.DataType, error) {
	data := make(mmdbtype.Map)
	for _, key := range val.MapKeys() {
		parseValue, err := convertToMMDBType(val.MapIndex(key))
		if err != nil {
			return nil, err
		}
		data[mmdbtype.String(key.String())] = parseValue
	}
	return data, nil
}

// handleSliceAndArray processes the provided reflect.Value of a slice or array type,
// and converts it into a MMDB slice type.
func handleSliceAndArray(val reflect.Value) (mmdbtype.DataType, error) {
	subData := make(mmdbtype.Slice, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		parseValue, err := convertToMMDBType(val.Index(i))
		if err != nil {
			return nil, err
		}
		subData = append(subData, parseValue)
	}
	return subData, nil
}

// handleSimpleType processes the provided reflect.Value of a simple type,
// and converts it into a MMDB data type.
func handleSimpleType(val reflect.Value) mmdbtype.DataType {
	switch val.Kind() {
	case reflect.String:
		return mmdbtype.String(val.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return mmdbtype.Int32(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16:
		return mmdbtype.Uint16(val.Uint())
	case reflect.Uint32:
		return mmdbtype.Uint32(val.Uint())
	case reflect.Uint64:
		return mmdbtype.Uint64(val.Uint())
	case reflect.Float32:
		return mmdbtype.Float32(val.Float())
	case reflect.Float64:
		return mmdbtype.Float64(val.Float())
	case reflect.Bool:
		return mmdbtype.Bool(val.Bool())
	}

	return mmdbtype.String("")
}
