/*
 * Copyright (c) 2022 shenjunzheng@gmail.com
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
	"bytes"
	"fmt"
	"sort"
)

// TextParser 文本解析器
type TextParser struct {
	Text           string
	IPv4FillResult func(str string) string
	IPv6FillResult func(str string) string

	segments map[int]Segment
}

// NewTextParser 初始化文本解析器
func NewTextParser(text string) *TextParser {
	return &TextParser{
		Text:     text,
		segments: make(map[int]Segment),
	}
}

// Parse 解析文本
func (p *TextParser) Parse() *TextParser {
	return p.ParseIPv4().ParseIPv6()
}

// ParseIPv4 解析文本中的 IPv4 片段
func (p *TextParser) ParseIPv4() *TextParser {
	ipv4Index := IPv4Regexp.FindAllStringIndex(p.Text, -1)
	if len(ipv4Index) > 0 {
		if p.segments == nil {
			p.segments = make(map[int]Segment)
		}
		for _, v := range ipv4Index {
			seg := Segment{
				Start:   v[0],
				End:     v[1],
				Type:    "ipv4",
				Content: p.Text[v[0]:v[1]],
			}
			if p.IPv4FillResult != nil {
				seg.Result = p.IPv4FillResult(seg.Content)
			}
			p.segments[v[0]] = seg
		}
	}
	return p
}

// ParseIPv6 解析文本中的 IPv6 片段
func (p *TextParser) ParseIPv6() *TextParser {
	ipv6Index := IPv6Regexp.FindAllStringIndex(p.Text, -1)
	if len(ipv6Index) > 0 {
		if p.segments == nil {
			p.segments = make(map[int]Segment)
		}
		for _, v := range ipv6Index {
			seg := Segment{
				Start:   v[0],
				End:     v[1],
				Type:    "ipv6",
				Content: p.Text[v[0]:v[1]],
			}
			if p.IPv6FillResult != nil {
				seg.Result = p.IPv6FillResult(seg.Content)
			}
			p.segments[v[0]] = seg
		}
	}
	return p
}

// String 返回解析后的文本
func (p *TextParser) String() string {
	if p.segments == nil {
		return p.Text
	}
	segments := make([]Segment, 0, len(p.segments))
	for _, segment := range p.segments {
		segments = append(segments, segment)
	}
	sort.Sort(SorterSegment(segments))

	var buffer bytes.Buffer
	var start int
	for _, v := range segments {
		if v.Start < start {
			continue
		}
		if v.Start > start {
			buffer.WriteString(p.Text[start:v.Start])
		}
		buffer.WriteString(v.Content)
		if len(v.Result) != 0 {
			buffer.WriteString(fmt.Sprintf(" [%s] ", v.Result))
		}
		start = v.End
	}
	buffer.WriteString(p.Text[start:])
	return buffer.String()
}

// Segment 文本片段
type Segment struct {
	Start   int
	End     int
	Type    string
	Content string
	Result  string
}

// SorterSegment 用于对 Segment 进行排序，按照 Start 字段进行排序
type SorterSegment []Segment

func (s SorterSegment) Len() int {
	return len(s)
}
func (s SorterSegment) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SorterSegment) Less(i, j int) bool {
	return s[i].Start < s[j].Start
}
