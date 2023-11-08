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

package domainlist

import (
	"bufio"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"

	"github.com/sjzar/ips/domainlist/data"
	"github.com/sjzar/ips/pkg/model"
)

const (
	DataSep = "\t"
)

// DomainList is a global map that holds domain information with the domain as the key.
var DomainList map[string]string

func init() {
	DomainList = make(map[string]string)

	r := strings.NewReader(strings.Join([]string{data.PlatformDomainList, data.ApplicationDomainList, data.OverseasDomainList}, "\n"))
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		split := strings.SplitN(string(line), DataSep, 2)
		if len(split) < 2 {
			continue
		}
		DomainList[split[0]] = split[1]
	}
}

// GetDomainInfo returns the domain information for the given domain.
// It attempts to find the main domain, then retrieves and parses its associated data.
// Returns a pointer to a DomainInfo object and true if successful, nil and false otherwise.
func GetDomainInfo(domain string) (*model.DomainInfo, bool) {
	mainDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return nil, false
	}

	info, ok := DomainList[mainDomain]
	if !ok {
		return nil, false
	}

	values, err := url.ParseQuery(info)
	if err != nil {
		return nil, false
	}

	ret := &model.DomainInfo{
		Domain:     domain,
		MainDomain: mainDomain,
		Data:       make(map[string]string),
	}
	for k, v := range values {
		ret.Data[k] = v[0]
	}

	return ret, ok
}

// GetDomainName returns the domain name for the given domain.
// It uses GetDomainInfo to retrieve the domain's information and extracts the "name" field.
// Returns the name and true if successful, an empty string and false otherwise.
func GetDomainName(domain string) (string, bool) {
	info, ok := GetDomainInfo(domain)
	if !ok {
		return "", false
	}

	name, ok := info.Data["name"]
	return name, ok
}
