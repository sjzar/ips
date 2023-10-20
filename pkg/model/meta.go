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

const (
	// IPv4 represents the version for IPv4 IP database.
	IPv4 = 1

	// IPv6 represents the version for IPv6 IP database.
	IPv6 = 2

	// MetaVersion represents the version of the metadata.
	MetaVersion = 1
)

// Meta represents the metadata for the IP database.
type Meta struct {

	// MetaVersion indicates the version of the metadata.
	MetaVersion int

	// Format represents the data format of the IP database.
	Format string

	// IPVersion indicates the version of the IP database (IPv4 or IPv6).
	IPVersion int

	// Fields lists the data fields available in the IP database.
	Fields []string

	// FieldAlias maps common field names to database-specific field names.
	FieldAlias map[string]string
}

// IsIPv4Support checks if the metadata supports IPv4.
func (m *Meta) IsIPv4Support() bool {
	return m.IPVersion&IPv4 != 0
}

// IsIPv6Support checks if the metadata supports IPv6.
func (m *Meta) IsIPv6Support() bool {
	return m.IPVersion&IPv6 != 0
}

// AddCommonFieldAlias adds aliases for common fields to the metadata.
// It checks if the provided database field exists in the metadata's fields.
func (m *Meta) AddCommonFieldAlias(fieldAlias map[string]string) {
	if m.FieldAlias == nil {
		m.FieldAlias = make(map[string]string)
	}

	fields := make(map[string]bool)
	for _, field := range m.Fields {
		fields[field] = true
	}

	for commonField, dbField := range fieldAlias {
		if _, ok := fields[dbField]; !ok {
			continue
		}
		m.FieldAlias[commonField] = dbField
	}
}

// SupportFields returns a map of fields that are supported by the metadata.
// It includes both the original fields and the alias fields.
func (m *Meta) SupportFields() map[string]bool {
	fields := make(map[string]bool)
	for _, field := range m.Fields {
		fields[field] = true
	}

	for commonField := range m.FieldAlias {
		fields[commonField] = true
	}

	return fields
}
