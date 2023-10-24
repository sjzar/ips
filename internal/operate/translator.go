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

package operate

import (
	"github.com/sjzar/ips/format/geo"
	"github.com/sjzar/ips/pkg/model"
)

// Translator is a structure that provides functionality to translate IPInfo's geolocation fields.
type Translator struct {
	TargetLang string // Target language for translation
}

// NewTranslator initializes a new Translator for a given target language.
func NewTranslator(targetLang string) (*Translator, error) {
	if len(targetLang) != 0 {
		if err := geo.SetLanguage(targetLang); err != nil {
			return nil, err
		}
	}
	return &Translator{
		TargetLang: targetLang,
	}, nil
}

// Do translates the geolocation fields in the provided IPInfo instance.\
func (t *Translator) Do(info *model.IPInfo) error {

	// List of geolocation fields to be translated
	fields := []string{model.Continent, model.Country, model.Province, model.City}

	// Iterate through each field and translate it
	for _, field := range fields {
		if value, ok := info.Data[field]; ok {
			info.Data[field] = geo.Translate(field, value)
		} else if alias, ok := info.FieldAlias[field]; ok {
			if value, ok := info.Data[alias]; ok {
				info.Data[alias] = geo.Translate(field, value)
			}
		}
	}

	return nil
}
