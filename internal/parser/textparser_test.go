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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextParser(t *testing.T) {
	ast := assert.New(t)

	type instance struct {
		str      string
		segments []Segment
	}

	instances := []instance{
		{str: "", segments: []Segment{}},
		{str: "123", segments: []Segment{{Start: 0, End: 3, Type: TextTypeText, Content: "123"}}},
		{str: "1.1.1.1", segments: []Segment{{Start: 0, End: 7, Type: TextTypeIPv4, Content: "1.1.1.1"}}},
		{str: "abc1.1.1.1def", segments: []Segment{
			{Start: 0, End: 3, Type: TextTypeText, Content: "abc"},
			{Start: 3, End: 10, Type: TextTypeIPv4, Content: "1.1.1.1"},
			{Start: 10, End: 13, Type: TextTypeText, Content: "def"}},
		},
		{str: "666.123.231.321g99::80z", segments: []Segment{
			{Start: 0, End: 1, Type: TextTypeText, Content: "6"},
			{Start: 1, End: 14, Type: TextTypeIPv4, Content: "66.123.231.32"},
			{Start: 14, End: 16, Type: TextTypeText, Content: "1g"},
			{Start: 16, End: 22, Type: TextTypeIPv6, Content: "99::80"},
			{Start: 22, End: 23, Type: TextTypeText, Content: "z"}},
		},
		{str: "::ffff:1.2.3.4.5.6.7.8", segments: []Segment{
			{Start: 0, End: 14, Type: TextTypeIPv6, Content: "::ffff:1.2.3.4"},
			{Start: 14, End: 15, Type: TextTypeText, Content: "."},
			{Start: 15, End: 22, Type: TextTypeIPv4, Content: "5.6.7.8"}},
		},
	}

	for index, inst := range instances {
		parser := NewTextParser(inst.str)
		parser.Parse()
		ast.Equal(inst.segments, parser.Segments, "index: %d str: %s", index, inst.str)
	}
}
