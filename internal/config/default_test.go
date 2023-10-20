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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDefault(t *testing.T) {
	ast := assert.New(t)

	type Inner struct {
		City string `default:"San Francisco"`
	}

	type Example struct {
		Int         int                         `default:"1"`
		String      string                      `default:"hello"`
		SliceString []string                    `default:"[\"hello\",\"world\"]"`
		Struct      Inner                       `default:"{\"City\":\"Shanghai\"}"`
		StructPtr   *Inner                      `default:"{\"City\":\"Los Angeles\"}"`
		StructNoTag Inner                       `default:"{\"City\":\"Los Angeles\"}"`
		Complex     map[string]map[string]Inner `default:"{\"location\":{\"city\":{\"City\":\"Los Angeles\"}}}"`
		Complex2    map[string][][]Inner        `default:"{\"location\":[[{\"City\":\"Los Angeles\"}]]}"`
		Interface   interface{}                 `default:"{\"City\":\"Los Angeles\"}"`
	}

	tests := []struct {
		input    Example
		expected Example
	}{
		{
			input: Example{},
			expected: Example{
				Int:         1,
				String:      "hello",
				SliceString: []string{"hello", "world"},
				Struct:      Inner{"Shanghai"},
				StructPtr:   &Inner{"Los Angeles"},
				StructNoTag: Inner{"Los Angeles"},
				Complex:     map[string]map[string]Inner{"location": {"city": {"Los Angeles"}}},
				Complex2:    map[string][][]Inner{"location": {{{"Los Angeles"}}}},
				Interface:   nil,
			},
		},
		{
			input: Example{
				Int:         9,
				String:      "world",
				SliceString: []string{"foo", "bar"},
				Struct:      Inner{"WenZhou"},
				StructPtr:   &Inner{"Shanghai"},
				StructNoTag: Inner{"Shanghai"},
				Complex:     map[string]map[string]Inner{"location": {"city": {"Shanghai"}}},
				Complex2:    map[string][][]Inner{"location": {{{"Shanghai"}}}},
				Interface:   &Inner{},
			},
			expected: Example{
				Int:         9,
				String:      "world",
				SliceString: []string{"foo", "bar"},
				Struct:      Inner{"WenZhou"},
				StructPtr:   &Inner{"Shanghai"},
				StructNoTag: Inner{"Shanghai"},
				Complex:     map[string]map[string]Inner{"location": {"city": {"Shanghai"}}},
				Complex2:    map[string][][]Inner{"location": {{{"Shanghai"}}}},
				Interface:   &Inner{"Los Angeles"},
			},
		},
	}

	for _, test := range tests {
		SetDefault(&test.input)
		ast.Equal(test.expected, test.input)
	}
}

func TestHandleSimpleType(t *testing.T) {
	ast := assert.New(t)

	type Example struct {
		Int    int     `default:"1"`
		Uint   uint    `default:"1"`
		Float  float64 `default:"1.1"`
		Bool   bool    `default:"true"`
		String string  `default:"hello"`
	}

	tests := []struct {
		input    Example
		expected Example
	}{
		{
			input: Example{},
			expected: Example{
				Int:    1,
				Uint:   1,
				Float:  1.1,
				Bool:   true,
				String: "hello",
			},
		},
		{
			input: Example{
				Int:    9,
				Uint:   10,
				Float:  11.1,
				String: "world",
			},
			expected: Example{
				Int:    9,
				Uint:   10,
				Float:  11.1,
				Bool:   true,
				String: "world",
			},
		},
	}

	for _, test := range tests {
		SetDefault(&test.input)
		ast.Equal(test.expected, test.input)
	}
}

func TestHandleSliceArray(t *testing.T) {
	ast := assert.New(t)

	type Inner struct {
		City string `default:"San Francisco"`
	}

	type Example struct {
		SliceInt    []int    `default:"[1,2,3]"`
		SliceString []string `default:"[\"hello\",\"world\"]"`
		SlicePtr    []*Inner `default:"[{\"City\":\"Los Angeles\"}]"`
		SliceNoTag  []Inner
		ArrayInt    [3]int    `default:"[1,2]"`
		ArrayString [2]string `default:"[\"hello\",\"world\",\"more\"]"`
		ArrayStruct [2]Inner  `default:"[{\"City\":\"Los Angeles\"}]"`
		ArrayPtr    [2]*Inner `default:"[{\"City\":\"Los Angeles\"}]"`
		ArrayNoTag  [2]Inner
	}

	tests := []struct {
		input    Example
		expected Example
	}{
		{
			input: Example{},
			expected: Example{
				SliceInt:    []int{1, 2, 3},
				SliceString: []string{"hello", "world"},
				SlicePtr:    []*Inner{{"Los Angeles"}},
				ArrayInt:    [3]int{1, 2, 0},
				ArrayString: [2]string{"hello", "world"},
				ArrayStruct: [2]Inner{{"Los Angeles"}, {"San Francisco"}},
				ArrayPtr:    [2]*Inner{{"Los Angeles"}, nil},
				ArrayNoTag:  [2]Inner{{"San Francisco"}, {"San Francisco"}},
			},
		},
		{
			input: Example{
				SliceInt:    []int{4, 5, 6},
				SliceString: []string{"foo", "bar"},
				SlicePtr:    []*Inner{{"Los Angeles"}},
				ArrayInt:    [3]int{4, 5, 6},
				ArrayString: [2]string{"hello", "world"},
				ArrayStruct: [2]Inner{{"WenZhou"}, {"Shanghai"}},
				ArrayPtr:    [2]*Inner{{"WenZhou"}, {"Shanghai"}},
				ArrayNoTag:  [2]Inner{{"WenZhou"}, {"Shanghai"}},
			},
			expected: Example{
				SliceInt:    []int{4, 5, 6},
				SliceString: []string{"foo", "bar"},
				SlicePtr:    []*Inner{{"Los Angeles"}},
				ArrayInt:    [3]int{4, 5, 6},
				ArrayString: [2]string{"hello", "world"},
				ArrayStruct: [2]Inner{{"WenZhou"}, {"Shanghai"}},
				ArrayPtr:    [2]*Inner{{"WenZhou"}, {"Shanghai"}},
				ArrayNoTag:  [2]Inner{{"WenZhou"}, {"Shanghai"}},
			},
		},
	}

	for _, test := range tests {
		SetDefault(&test.input)
		ast.Equal(test.expected, test.input)
	}
}

func TestHandleMap(t *testing.T) {
	ast := assert.New(t)

	type Inner struct {
		City string `default:"San Francisco"`
	}

	type Example struct {
		MapInt               map[string]int      `default:"{\"country\":1,\"state\":2}"`
		MapString            map[string]string   `default:"{\"country\":\"US\",\"state\":\"CA\"}"`
		MapStruct            map[string]Inner    `default:"{\"location\":{\"City\":\"Los Angeles\"}}"`
		MapPtr               map[string]*Inner   `default:"{\"location\":{\"City\":\"Los Angeles\"}}"`
		MapSlice             map[string][]string `default:"{\"location\":[\"Los Angeles\"]}"`
		MapPtrWithoutDefault map[string]*Inner
	}

	tests := []struct {
		input    Example
		expected Example
	}{
		{
			input: Example{},
			expected: Example{
				MapInt:    map[string]int{"country": 1, "state": 2},
				MapString: map[string]string{"country": "US", "state": "CA"},
				MapStruct: map[string]Inner{"location": {"Los Angeles"}},
				MapPtr:    map[string]*Inner{"location": {"Los Angeles"}},
				MapSlice:  map[string][]string{"location": {"Los Angeles"}},
			},
		},
	}

	for _, test := range tests {
		SetDefault(&test.input)
		ast.Equal(test.expected, test.input)
	}

}
