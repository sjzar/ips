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
	"fmt"
	"strings"

	"github.com/sjzar/ips/internal/util"
)

// AlfredItem represents an item to be displayed in Alfred's result list.
type AlfredItem struct {
	Title    string     `json:"title"`
	Subtitle string     `json:"subtitle"`
	Arg      string     `json:"arg"`
	Icon     AlfredIcon `json:"icon"`
	Valid    bool       `json:"valid"`
	Text     AlfredText `json:"text"`
}

// AlfredIcon represents the icon for an AlfredItem.
type AlfredIcon struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

// AlfredText provides additional text information for an AlfredItem.
type AlfredText struct {
	Copy string `json:"copy"`
}

// DataList holds a list of items to be displayed in Alfred's result list.
type DataList struct {
	Items []interface{} `json:"items"`
}

// AddItem appends a new item to the DataList's Items slice.
func (d *DataList) AddItem(item interface{}) {
	if d.Items == nil {
		d.Items = make([]interface{}, 0)
	}
	d.Items = append(d.Items, item)
}

// AddAlfredItemByIPInfo creates an AlfredItem from the provided IPInfo
// and adds it to the DataList.
func (d *DataList) AddAlfredItemByIPInfo(info *IPInfo) {
	values := strings.Join(util.DeleteEmptyValue(info.Values()), " ")
	item := AlfredItem{
		Title:    fmt.Sprintf("%s [%s]", info.IP, values),
		Subtitle: "Copy to clipboard",
		Arg:      values,
		Icon:     AlfredIcon{},
		Valid:    true,
		Text: AlfredText{
			Copy: values,
		},
	}
	d.AddItem(item)
}

// AddAlfredItemByDomainInfo creates an AlfredItem from the provided DomainInfo
// and adds it to the DataList.
func (d *DataList) AddAlfredItemByDomainInfo(info *DomainInfo) {
	values := strings.Join(util.DeleteEmptyValue(info.Values()), " ")
	item := AlfredItem{
		Title:    fmt.Sprintf("%s [%s]", info.Domain, values),
		Subtitle: "Copy to clipboard",
		Arg:      values,
		Icon:     AlfredIcon{},
		Valid:    true,
		Text: AlfredText{
			Copy: values,
		},
	}
	d.AddItem(item)
}

// AddAlfredItemEmpty adds a default "Not found" AlfredItem to the DataList
// if the list is empty.
func (d *DataList) AddAlfredItemEmpty() {
	if len(d.Items) > 0 {
		return
	}
	item := AlfredItem{
		Title:    "Not found",
		Subtitle: "No information found",
		Arg:      "",
		Icon:     AlfredIcon{},
		Valid:    false,
		Text:     AlfredText{},
	}
	d.AddItem(item)
}
