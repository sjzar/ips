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

package config

import (
	"encoding/json"
	"reflect"
	"strconv"
)

// DefaultTag is the default tag used to identify default values for struct fields.
var DefaultTag string

func init() {
	DefaultTag = "default"
}

// SetDefaultTag updates the tag used to identify default values.
func SetDefaultTag(tag string) {
	DefaultTag = tag
}

// SetDefault sets the default values for a given interface{}.
// The interface{} must be a pointer to a struct, bcs the default values are SET on the struct fields.
func SetDefault(v interface{}) {
	if v == nil {
		return
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	setDefault(val, "")
}

// setDefault recursively sets default values based on the struct's tags.
func setDefault(val reflect.Value, tag string) {

	if !val.CanSet() {
		return
	}

	switch val.Kind() {
	case reflect.Struct:
		handleStruct(val, tag)
	case reflect.Ptr:
		handlePtr(val, tag)
	case reflect.Interface:
		handleInterface(val, tag)
	case reflect.Map:
		handleMap(val, tag)
	case reflect.Slice, reflect.Array:
		handleSliceArray(val, tag)
	default:
		handleSimpleType(val, tag)
	}
}

// handleSimpleType handles the assignment of default values for simple data types.
func handleSimpleType(val reflect.Value, tag string) {

	if len(tag) == 0 || !val.IsZero() || !val.CanSet() {
		return
	}

	// Assign appropriate default values based on the data type.
	switch val.Kind() {
	case reflect.String:
		val.SetString(tag)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(tag, 10, 64); err == nil {
			val.SetInt(intValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintVal, err := strconv.ParseUint(tag, 10, 64); err == nil {
			val.SetUint(uintVal)
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(tag, 64); err == nil {
			val.SetFloat(floatValue)
		}
	case reflect.Bool:
		if boolVal, err := strconv.ParseBool(tag); err == nil {
			val.SetBool(boolVal)
		}
	}
}

// handleStruct processes struct type fields and sets default values based on their tags.
func handleStruct(val reflect.Value, tag string) {

	if val.Kind() != reflect.Struct {
		return
	}

	if !val.IsZero() || len(tag) == 0 {
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)
			setDefault(field, fieldType.Tag.Get(DefaultTag))
		}
		return
	}

	if !val.CanSet() {
		return
	}

	// If tag is provided, unmarshal the default value and set it.
	newStruct := reflect.New(val.Type()).Interface()
	if err := json.Unmarshal([]byte(tag), newStruct); err == nil {
		newStructVal := reflect.ValueOf(newStruct).Elem()
		for i := 0; i < newStructVal.NumField(); i++ {
			field := newStructVal.Field(i)
			fieldType := newStructVal.Type().Field(i)
			setDefault(field, fieldType.Tag.Get(DefaultTag))
		}
		val.Set(newStructVal)
	}
}

// handleSliceArray sets default values for slice or array types.
func handleSliceArray(val reflect.Value, tag string) {

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return
	}

	if !val.IsZero() || len(tag) == 0 {
		for i := 0; i < val.Len(); i++ {
			setDefault(val.Index(i), "")
		}
		return
	}

	if !val.CanSet() {
		return
	}

	// array 提前初始化子结构体，解决 default 长度小于 array 长度的问题
	if val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			setDefault(val.Index(i), "")
		}
	}

	// If tag is provided, unmarshal the default value and set it.
	newSlice := reflect.New(reflect.SliceOf(val.Type().Elem())).Interface()
	if err := json.Unmarshal([]byte(tag), newSlice); err == nil {
		sliceValue := reflect.ValueOf(newSlice).Elem()
		for j := 0; j < sliceValue.Len(); j++ {
			v := sliceValue.Index(j)
			setDefault(v, "")
			if val.Kind() == reflect.Array {
				if j >= val.Len() {
					return
				}
				val.Index(j).Set(v)
			} else {
				val.Set(reflect.Append(val, v))
			}
		}
	}
}

// handleMap sets default values for map types.
func handleMap(val reflect.Value, tag string) {

	if val.Kind() != reflect.Map {
		return
	}

	if !val.IsZero() || len(tag) == 0 {
		for _, key := range val.MapKeys() {
			setDefault(val.MapIndex(key), "")
		}
		return
	}

	if !val.CanSet() {
		return
	}

	// If tag is provided, unmarshal the default value and set it.
	newMap := reflect.New(val.Type()).Interface()
	if err := json.Unmarshal([]byte(tag), newMap); err == nil {
		newMapVal := reflect.ValueOf(newMap).Elem()
		for _, k := range newMapVal.MapKeys() {
			v := newMapVal.MapIndex(k)
			setDefault(v, "")
		}
		val.Set(newMapVal)
	}
}

// handlePtr handles pointer types and sets default values based on their tags.
func handlePtr(val reflect.Value, tag string) {

	if val.Kind() != reflect.Ptr {
		return
	}

	if !val.IsZero() || len(tag) == 0 || !val.CanSet() {
		return
	}

	// If tag is provided, unmarshal the default value and set it.
	newPtr := reflect.New(val.Type()).Interface()
	if err := json.Unmarshal([]byte(tag), newPtr); err == nil {
		newPtrVal := reflect.ValueOf(newPtr).Elem()
		setDefault(newPtrVal, "")
		val.Set(newPtrVal)
	}
}

// handleInterface processes interface types and sets default values based on their tags.
// support pointer interface only, bcs the default values are SET on the struct fields.
func handleInterface(val reflect.Value, tag string) {

	if val.Kind() != reflect.Interface {
		return
	}

	if val.IsNil() {
		return
	}

	if val.Elem().Kind() != reflect.Ptr {
		return
	}

	// Recursively set the default for the inner value of the interface.
	setDefault(reflect.ValueOf(val.Elem().Interface()).Elem(), tag)
}
