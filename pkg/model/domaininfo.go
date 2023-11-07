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
	"sort"
)

// DomainInfo holds details about a specific domain, including its name, main domain, and associated data.
type DomainInfo struct {

	// Domain is the domain name.
	Domain string `json:"domain"`

	// MainDomain is the main domain name.
	MainDomain string `json:"main_domain"`

	// Data holds the actual information related to the domain.
	Data map[string]string `json:"data"`
}

// Values extracts the values from the Data map, sorts them alphabetically, and returns them as a slice of strings.
// If the DomainInfo receiver or its Data map is nil, it returns nil to indicate the absence of data.
func (d *DomainInfo) Values() []string {
	if d == nil || d.Data == nil {
		return nil
	}
	keys := make([]string, 0, len(d.Data))
	for k := range d.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := make([]string, 0, len(d.Data))
	for _, v := range keys {
		values = append(values, d.Data[v])
	}

	return values
}
