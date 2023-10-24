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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertMapToFields(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		name            string
		input           map[string]interface{}
		enableFullField bool
		expected        map[string]string
		expectError     bool
	}

	testCases := []testCase{
		{
			name: "Handle String Type",
			input: map[string]interface{}{
				"simpleString": "Hello, World!",
			},
			enableFullField: false,
			expected: map[string]string{
				"simpleString": "Hello, World!",
			},
			expectError: false,
		},
		{
			name: "Handle Integer Types",
			input: map[string]interface{}{
				"simpleInt": 123,
				"smallInt":  int8(8),
			},
			enableFullField: false,
			expected: map[string]string{
				"simpleInt": "123",
				"smallInt":  "8",
			},
			expectError: false,
		},
		{
			name: "Handle Unsigned Integer Types",
			input: map[string]interface{}{
				"simpleUint": uint(123),
				"smallUint":  uint8(8),
			},
			enableFullField: false,
			expected: map[string]string{
				"simpleUint": "123",
				"smallUint":  "8",
			},
			expectError: false,
		},
		{
			name: "Handle Float Types",
			input: map[string]interface{}{
				"simpleFloat": 123.456,
				"smallFloat":  float32(8.9),
			},
			enableFullField: false,
			expected: map[string]string{
				"simpleFloat": "123.456000",
				"smallFloat":  "8.900000",
			},
			expectError: false,
		},
		{
			name: "Handle Boolean Type",
			input: map[string]interface{}{
				"boolValue": true,
			},
			enableFullField: false,
			expected: map[string]string{
				"boolValue": "true",
			},
			expectError: false,
		},
		{
			name: "Handle Complex Types",
			input: map[string]interface{}{
				"complexStruct": struct {
					Field1 string
					Field2 int
				}{
					Field1: "Hello",
					Field2: 42,
				},
				"complexMap": map[string]interface{}{
					"subKey1": "value1",
					"subKey2": 1234,
				},
				"complexArray": []interface{}{"elem1", "elem2", 5678},
			},
			enableFullField: false,
			expected: map[string]string{
				"complexStruct_Field1": "Hello",
				"complexStruct_Field2": "42",
				"complexMap_subKey1":   "value1",
				"complexMap_subKey2":   "1234",
				"complexArray":         "elem1,elem2,5678",
			},
			expectError: false,
		},
		{
			name: "Handle String Type with Full Field",
			input: map[string]interface{}{
				"simpleString": "Hello, World!",
			},
			enableFullField: true,
			expected: map[string]string{
				"simpleString": "\"Hello, World!\"",
			},
			expectError: false,
		},
		{
			name: "Handle Integer Types with Full Field",
			input: map[string]interface{}{
				"simpleInt": 123,
				"smallInt":  int8(8),
			},
			enableFullField: true,
			expected: map[string]string{
				"simpleInt": "123",
				"smallInt":  "8",
			},
			expectError: false,
		},
		{
			name: "Handle Complex Types with Full Field",
			input: map[string]interface{}{
				"complexStruct": struct {
					Field1 string
					Field2 int
				}{
					Field1: "Hello",
					Field2: 42,
				},
				"complexMap": map[string]interface{}{
					"subKey1": "value1",
					"subKey2": 1234,
				},
				"complexArray": []interface{}{"elem1", "elem2", 5678},
			},
			enableFullField: true,
			expected: map[string]string{
				"complexStruct": "{\"Field1\":\"Hello\",\"Field2\":42}",
				"complexMap":    "{\"subKey1\":\"value1\",\"subKey2\":1234}",
				"complexArray":  "[\"elem1\",\"elem2\",5678]",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		output, err := ConvertMapToFields(tc.input, tc.enableFullField)
		if tc.expectError {
			ast.Error(err)
		} else {
			ast.NoError(err)
			ast.Equal(tc.expected, output)
		}
	}
}
