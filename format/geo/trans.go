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
	"github.com/sjzar/ips/pkg/errors"
)

var DatabaseLanguage string

// init sets the default database language upon initialization.
func init() {
	DatabaseLanguage = DefaultLanguage
}

// SetDBLanguage sets the language for the database.
// It returns an error if the provided language is not supported.
func SetDBLanguage(lang string) error {
	for _, l := range SupportedLanguages {
		if l == lang {
			DatabaseLanguage = lang
			return nil
		}
	}

	return errors.ErrUnsupportedLanguage
}

// Translate translates the provided text from the database language to the application's current language.
// If the text cannot be translated, it returns the original text.
func Translate(field, text string) string {
	return translate(DatabaseLanguage, Language, field, text)
}

// translate translates the provided text from the source language to the target language.
// If the text cannot be translated, it returns the original text.
func translate(sourceLang, targetLang, field, text string) string {
	// If source and target languages are the same, return the original text
	if sourceLang == targetLang {
		return text
	}

	// Try to get the geolocation info for the text in the source language
	nameInfos := GetNameInfos(field, sourceLang)
	str, ok := nameInfos[text]
	if !ok {
		return text
	}
	info, ok := ParseGeoInfo(str)
	if !ok {
		return text
	}

	// Return the translated name in the target language
	return info.Name(targetLang)
}
