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
	"sort"
)

const (
	// TextTypeIPv4 represents the IPv4 type in the text.
	TextTypeIPv4 = "ipv4"

	// TextTypeIPv6 represents the IPv6 type in the text.
	TextTypeIPv6 = "ipv6"

	// TextTypeDomain represents the domain type in the text.
	TextTypeDomain = "domain"

	// TextTypeText represents any other text type.
	TextTypeText = "text"
)

// TextParser is a parser designed to extract IPv4, IPv6, and other text segments from a given text.
type TextParser struct {
	Text     string    // Text to be parsed
	Segments []Segment // Parsed segments from the text
}

// NewTextParser initializes and returns a new TextParser.
func NewTextParser(text string) *TextParser {
	return &TextParser{
		Text:     text,
		Segments: make([]Segment, 0),
	}
}

// Parse parses the text to extract IPv4, IPv6, and other segments.
func (t *TextParser) Parse() *TextParser {
	return t.ParseIPv4().ParseIPv6().Distinct()
}

// ParseIPv4 extracts IPv4 segments from the text.
func (t *TextParser) ParseIPv4() *TextParser {
	index := IPv4Regexp.FindAllStringIndex(t.Text, -1)
	if len(index) > 0 {
		if t.Segments == nil {
			t.Segments = make([]Segment, 0, len(index))
		}
		for _, v := range index {
			seg := Segment{
				Start:   v[0],
				End:     v[1],
				Type:    TextTypeIPv4,
				Content: t.Text[v[0]:v[1]],
			}
			t.Segments = append(t.Segments, seg)
		}
	}
	return t
}

// ParseIPv6 extracts IPv6 segments from the text.
func (t *TextParser) ParseIPv6() *TextParser {
	index := IPv6Regexp.FindAllStringIndex(t.Text, -1)
	if len(index) > 0 {
		if t.Segments == nil {
			t.Segments = make([]Segment, 0, len(index))
		}
		for _, v := range index {
			seg := Segment{
				Start:   v[0],
				End:     v[1],
				Type:    TextTypeIPv6,
				Content: t.Text[v[0]:v[1]],
			}
			t.Segments = append(t.Segments, seg)
		}
	}
	return t
}

// Distinct removes duplicate segments and completes the text slices.
func (t *TextParser) Distinct() *TextParser {
	sort.Sort(SorterSegment(t.Segments))
	segments := make([]Segment, 0, len(t.Segments)*2+1)
	var start int
	for _, v := range t.Segments {
		if v.Start < start {
			continue
		}
		if v.Start > start {
			seg := Segment{
				Start:   start,
				End:     v.Start,
				Type:    TextTypeText,
				Content: t.Text[start:v.Start],
			}
			segments = append(segments, seg)
		}
		segments = append(segments, v)
		start = v.End
	}
	if start < len(t.Text) {
		seg := Segment{
			Start:   start,
			End:     len(t.Text),
			Type:    TextTypeText,
			Content: t.Text[start:],
		}
		segments = append(segments, seg)
	}

	t.Segments = segments
	return t
}

// Segment represents a section or segment of the text.
type Segment struct {
	Start   int    // Start position of the segment
	End     int    // End position of the segment
	Type    string // Type of the segment (IPv4, IPv6, Domain, or Text)
	Content string // Actual content of the segment
}

// SorterSegment is a type to help in sorting a slice of Segments based on their Start position.
type SorterSegment []Segment

func (s SorterSegment) Len() int { return len(s) }

func (s SorterSegment) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SorterSegment) Less(i, j int) bool {
	if s[i].Start == s[j].Start {
		return TypeWeight(s[i].Type) > TypeWeight(s[j].Type)
	}
	return s[i].Start < s[j].Start
}

// TypeWeight provides a weight for segment types to help in sorting.
// IPv4 segments have the highest weight, followed by IPv6, then Domain, and then any other type.
func TypeWeight(t string) int {
	switch t {
	case TextTypeIPv4:
		return 3
	case TextTypeIPv6:
		return 2
	case TextTypeDomain:
		return 1
	default:
		return 0
	}
}
