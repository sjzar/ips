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

package geo

import (
	"bufio"
	"strings"

	"github.com/sjzar/ips/format/geo/data"
	"github.com/sjzar/ips/pkg/errors"
)

// Supported languages for geolocation information.
const (
	LangEnglish     = "en"
	LangChinese     = "zh-CN"
	LangRussian     = "ru"
	LangJapanese    = "ja"
	LangGerman      = "de"
	LangFrench      = "fr"
	LangSpanish     = "es"
	LangPortuguese  = "pt-BR"
	LangPersian     = "fa"
	LangKorean      = "ko"
	DefaultLanguage = LangChinese
)

var SupportedLanguages = []string{LangEnglish, LangChinese, LangRussian, LangJapanese, LangGerman,
	LangFrench, LangSpanish, LangPortuguese, LangPersian, LangKorean}

// Language to be used for names.
var Language string

func init() {
	Language = DefaultLanguage
}

// SetLanguage sets the global language for geolocation information.
func SetLanguage(lang string) error {
	for _, l := range SupportedLanguages {
		if l == lang {
			Language = lang
			return nil
		}
	}

	return errors.ErrUnsupportedLanguage
}

// IDInfos contains mapping from GeoNameID to its respective information.
var IDInfos map[string]string

// NameInfos contains a multilevel mapping from field -> language -> name -> information.
var NameInfos map[string]map[string]map[string]string

// LoadData populates geolocation data based on a given language and data sources.
func LoadData(lang string, data ...string) map[string]string {
	ret := make(map[string]string)
	for _, d := range data {
		scanner := bufio.NewScanner(strings.NewReader(d))
		for scanner.Scan() {
			line := scanner.Text()
			split := strings.SplitN(line, "\t", 3)
			if len(split) != 3 {
				continue
			}
			if len(lang) != 0 {
				if strings.Contains(split[1], lang) {
					split2 := strings.Split(split[1], "|")
					for _, l := range split2 {
						if strings.HasPrefix(l, lang) {
							ret[strings.TrimPrefix(l, lang+":")] = line
							break
						}
					}
				}
			} else {
				ret[split[0]] = line
			}
		}

	}
	return ret
}

// GetNameInfos retrieves geolocation information based on a given field and language.
func GetNameInfos(field string, lang string) map[string]string {
	if NameInfos == nil {
		NameInfos = make(map[string]map[string]map[string]string)
	}
	if NameInfos[field] == nil {
		NameInfos[field] = make(map[string]map[string]string)
	}
	if NameInfos[field][lang] == nil {
		switch field {
		case "continent":
			NameInfos[field][lang] = LoadData(lang, data.Continent)
		case "country", "country_name":
			NameInfos[field][lang] = LoadData(lang, data.Country)
		case "region", "province", "subdivisions":
			NameInfos[field][lang] = LoadData(lang, data.Region)
		case "city":
			NameInfos[field][lang] = LoadData(lang, data.City)
		default:
			NameInfos[field][lang] = LoadData(lang, data.Continent, data.Country, data.Region, data.City)
		}
	}

	return NameInfos[field][lang]
}
