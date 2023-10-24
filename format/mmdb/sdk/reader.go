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

package sdk

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/oschwald/maxminddb-golang"

	"github.com/sjzar/ips/format/geo"
	"github.com/sjzar/ips/ipnet"
	"github.com/sjzar/ips/pkg/model"
)

// Reader wraps the maxminddb Reader and provides additional functionalities
// for reading and parsing IP related data.
type Reader struct {
	db               *maxminddb.Reader // Database reader instance
	IPVersion        int               // IP versions supported by the database (IPv4, IPv6, or both)
	Fields           []string          // Fields present in the database
	DisableExtraData bool              // Whether to disable the use of extra data matched by GeoNameID
	UseFullField     bool              // Whether to use full field. If enabled, all data is combined into a single field as a JSON string.
}

// NewReader initializes a new Reader for the given MMDB file.
func NewReader(file string) (*Reader, error) {
	db, err := maxminddb.Open(file)
	if err != nil {
		return nil, err
	}

	ipVersion := 0
	switch db.Metadata.IPVersion {
	case 4:
		ipVersion |= model.IPv4
	case 6:
		ipVersion |= model.IPv4
		ipVersion |= model.IPv6
	}

	var m map[string]interface{}
	_, _, err = db.LookupNetwork(net.ParseIP("61.144.235.160"), &m)
	if err != nil {
		return nil, err
	}
	data, err := ConvertMapToFields(m, false)
	if err != nil {
		return nil, err
	}

	fields := make([]string, 0, len(data))
	for key := range data {
		fields = append(fields, key)
	}

	return &Reader{
		db:        db,
		IPVersion: ipVersion,
		Fields:    fields,
	}, nil
}

// Find looks up the given IP in the database and returns its associated range and data.
func (r *Reader) Find(ip net.IP) (*ipnet.Range, map[string]string, error) {
	var m map[string]interface{}
	ipNet, _, err := r.db.LookupNetwork(ip, &m)
	if err != nil {
		return nil, nil, err
	}

	data, err := ConvertMapToFields(m, r.UseFullField)
	if err != nil {
		return nil, nil, err
	}

	return ipnet.NewRange(ipNet), data, nil
}

// Close closes the underlying maxminddb Reader.
func (r *Reader) Close() error {
	return r.db.Close()
}

// ConvertMapToFields converts a map of data from MMDB into a map of strings,
// optionally using full fields (combining all data into a JSON string).
func ConvertMapToFields(m map[string]interface{}, useFullField bool) (map[string]string, error) {
	data := make(map[string]string)

	for key, _value := range m {
		if useFullField {
			jsonValue, err := json.Marshal(_value)
			if err != nil {
				return nil, err
			}
			data[key] = string(jsonValue)
			continue
		}

		value, err := ParseReflectValue(_value)
		if err != nil {
			return nil, err
		}
		switch v := value.(type) {
		case string:
			data[key] = v
		case map[string]string:
			for subKey, subValue := range v {
				data[key+"_"+subKey] = subValue
			}
		}
	}

	return DoPreProcessFieldAlias(data), nil
}

// ParseReflectValue processes an interface value and returns its converted string or map representation.
func ParseReflectValue(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return processReflectValue(val)
}

// processReflectValue processes a reflect.Value and returns its converted string or map representation.
func processReflectValue(val reflect.Value) (interface{}, error) {
	switch val.Kind() {
	case reflect.Struct:
		return handleStruct(val)
	case reflect.Ptr:
		return processReflectValue(val.Elem())
	case reflect.Interface:
		return processReflectValue(val.Elem())
	case reflect.Map:
		return handleMap(val)
	case reflect.Slice, reflect.Array:
		return handleSliceAndArray(val)
	default:
		return handleSimpleType(val), nil
	}
}

// handleStruct processes the provided reflect.Value of a struct type and converts it into a map of strings.
func handleStruct(val reflect.Value) (map[string]string, error) {
	data := make(map[string]string)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		parseValue, err := processReflectValue(field)
		if err != nil {
			return nil, err
		}
		switch v := parseValue.(type) {
		case string:
			data[fieldType.Name] = v
		case map[string]string:
			for subKey, subValue := range v {
				data[fieldType.Name+"_"+subKey] = subValue
			}
		}
	}
	return data, nil
}

// handleMap processes the provided reflect.Value of a map type and converts it into a map of strings or a single string.
func handleMap(val reflect.Value) (interface{}, error) {
	geoInfo, ok := geo.ParseInfoFromMMDB(val.Interface().(map[string]interface{}), false)
	if ok {
		return geoInfo.Name(geo.Language), nil
	}
	data := make(map[string]string)
	for _, key := range val.MapKeys() {
		parseValue, err := processReflectValue(val.MapIndex(key))
		if err != nil {
			return nil, err
		}
		switch v := parseValue.(type) {
		case string:
			data[key.String()] = v
		case map[string]string:
			for subKey, subValue := range v {
				data[key.String()+"_"+subKey] = subValue
			}
		}
	}
	return data, nil
}

// handleSliceAndArray processes the provided reflect.Value of a slice or array type and converts it into a string.
func handleSliceAndArray(val reflect.Value) (interface{}, error) {
	subData := make([]string, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		parseValue, err := processReflectValue(val.Index(i))
		if err != nil {
			return nil, err
		}
		switch v := parseValue.(type) {
		case string:
			subData = append(subData, v)
		case map[string]string:
			outJson, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			subData = append(subData, string(outJson))
		}
	}
	return strings.Join(subData, ","), nil
}

// handleSimpleType converts the provided reflect.Value of a basic type into a string.
func handleSimpleType(val reflect.Value) string {
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.6f", val.Float())
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	}

	return ""
}
