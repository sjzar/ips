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

package ips

import (
	"fmt"
	"strings"
)

func init() {
	// Set Persistence Default Values
	// viper.SetDefault("ipv4_file", []string{"qqwry.dat"})
}

const (
	// OutputTypeText represents the text output format.
	OutputTypeText = "text"

	// OutputTypeJSON represents the JSON output format.
	OutputTypeJSON = "json"

	// OutputTypeAlfred represents the Alfred output format.
	OutputTypeAlfred = "alfred"

	// DefaultFields represents the default output fields.
	DefaultFields = "country,province,city,isp"
)

// Config represents the application's configuration.
type Config struct {

	// Common
	// IPSDir specifies the working directory for IPS.
	IPSDir string `mapstructure:"-"`

	// Lang specifies the language for the output.
	Lang string `mapstructure:"lang"`

	// Find
	// IPv4File specifies the file is IPv4 database.
	IPv4File []string `mapstructure:"ipv4_file" default:"[\"qqwry.dat\"]"`

	// IPv4Format specifies the format for IPv4 database.
	IPv4Format []string `mapstructure:"ipv4_format"`

	// IPv6File specifies the file is IPv6 database.
	IPv6File []string `mapstructure:"ipv6_file" default:"[\"zxipv6wry.db\"]"`

	// IPv6Format specifies the format for IPv6 database.
	IPv6Format []string `mapstructure:"ipv6_format"`

	// HybridMode specifies the operational mode of the HybridReader.
	// It determines how the HybridReader processes and combines data from multiple IP database readers.
	// Accepted values are "comparison" and "aggregation":
	// - "comparison": Used for comparing data across different databases, where the output includes data from all readers.
	// - "aggregation": Used for creating a unified view of data by aggregating and supplementing missing fields from multiple databases. (default)
	HybridMode string `mapstructure:"hybrid_mode"`

	// Fields lists the output fields.
	// default is country, province, city, isp
	Fields string `mapstructure:"fields" default:"country,province,city,isp"`

	// UseDBFields indicates whether to use database fields. (default is common fields)
	UseDBFields bool `mapstructure:"use_db_fields"`

	// RewriteFiles lists the files for data rewriting.
	RewriteFiles string `mapstructure:"rewrite_files"`

	// OutputType specifies the type of the output. (default is text)
	OutputType string `mapstructure:"output_type"`

	// TextFormat specifies the format for text output.
	// It supports %origin and %values parameters.
	TextFormat string `mapstructure:"text_format" default:"%origin [%values] "`

	// TextValuesSep specifies the separator for values in text output. (default is space)
	TextValuesSep string `mapstructure:"text_values_sep" default:" "`

	// JsonIndent indicates whether the JSON output should be indented.
	JsonIndent bool `mapstructure:"json_indent"`

	// Dump & Pack
	// DPFields lists the output fields for dump and pack operations.
	// default is empty, means all fields
	DPFields string `mapstructure:"dp_fields"`

	// DPRewriterFiles lists the files for rewriting during dump and pack operations.
	DPRewriterFiles string `mapstructure:"dp_rewriter_files"`

	// Database
	// ReaderOption specifies the options for the reader.
	ReaderOption string `mapstructure:"reader_option"`

	// WriterOption specifies the options for the writer.
	WriterOption string `mapstructure:"writer_option"`

	// MyIP
	// LocalAddr specifies the local address (in IP format) that should be used for outbound connections.
	// Useful in systems with multiple network interfaces.
	LocalAddr string `mapstructure:"local_addr"`

	// MyIPCount defines the minimum number of detectors that should return the same IP
	// for the IP to be considered as the system's public IP.
	MyIPCount int `mapstructure:"myip_count" default:"3"`

	// MyIPTimeoutS specifies the maximum duration (in seconds) to wait for the detectors to return an IP.
	MyIPTimeoutS int `mapstructure:"myip_timeout_s"`

	// Service
	// Addr specifies the address for the service.
	Addr string `mapstructure:"addr" default:":6860"`
}

func (c *Config) ShowConfig(allKeys bool) string {
	str := fmt.Sprintf("ips dir:\t\t[%s]\n", c.IPSDir)
	if allKeys || len(c.Lang) > 0 {
		str += fmt.Sprintf("lang:\t\t\t[%s]\n", c.Lang)
	}
	str += fmt.Sprintf("ipv4_file(ipv4):\t[%s]\n", strings.Join(c.IPv4File, ","))
	if allKeys || len(c.IPv4Format) > 0 {
		str += fmt.Sprintf("ipv4_format:\t\t[%s]\n", strings.Join(c.IPv4Format, ","))
	}
	str += fmt.Sprintf("ipv6_file(ipv6):\t[%s]\n", strings.Join(c.IPv6File, ","))
	if allKeys || len(c.IPv6Format) > 0 {
		str += fmt.Sprintf("ipv6_format:\t\t[%s]\n", strings.Join(c.IPv6Format, ","))
	}
	if allKeys || len(c.HybridMode) > 0 {
		str += fmt.Sprintf("hybrid_mode:\t\t[%s]\n", c.HybridMode)
	}
	if allKeys || len(c.Fields) > 0 {
		str += fmt.Sprintf("fields:\t\t\t[%s]\n", c.Fields)
	}
	if allKeys || c.UseDBFields {
		str += fmt.Sprintf("use_db_fields:\t\t[%v]\n", c.UseDBFields)
	}
	if allKeys || len(c.RewriteFiles) > 0 {
		str += fmt.Sprintf("rewrite_files:\t\t[%s]\n", c.RewriteFiles)
	}
	if allKeys || len(c.OutputType) > 0 {
		str += fmt.Sprintf("output_type:\t\t[%s]\n", c.OutputType)
	}
	if allKeys || len(c.TextFormat) > 0 {
		str += fmt.Sprintf("text_format:\t\t[%s]\n", c.TextFormat)
	}
	if allKeys || len(c.TextValuesSep) > 0 {
		str += fmt.Sprintf("text_values_sep:\t[%s]\n", c.TextValuesSep)
	}
	if allKeys || c.JsonIndent {
		str += fmt.Sprintf("json_indent:\t\t[%v]\n", c.JsonIndent)
	}
	if allKeys || len(c.DPFields) > 0 {
		str += fmt.Sprintf("dp_fields:\t\t[%s]\n", c.DPFields)
	}
	if allKeys || len(c.DPRewriterFiles) > 0 {
		str += fmt.Sprintf("dp_rewriter_files:\t[%s]\n", c.DPRewriterFiles)
	}
	if allKeys || len(c.ReaderOption) > 0 {
		str += fmt.Sprintf("reader_option:\t\t[%s]\n", c.ReaderOption)
	}
	if allKeys || len(c.WriterOption) > 0 {
		str += fmt.Sprintf("writer_option:\t\t[%s]\n", c.WriterOption)
	}
	if allKeys || len(c.LocalAddr) > 0 {
		str += fmt.Sprintf("local_addr:\t\t[%s]\n", c.LocalAddr)
	}
	if allKeys || c.MyIPCount > 0 {
		str += fmt.Sprintf("myip_count:\t\t[%d]\n", c.MyIPCount)
	}
	if allKeys || c.MyIPTimeoutS > 0 {
		str += fmt.Sprintf("myip_timeout_s:\t\t[%d]\n", c.MyIPTimeoutS)
	}
	if allKeys || len(c.Addr) > 0 {
		str += fmt.Sprintf("addr:\t\t\t[%s]\n", c.Addr)
	}

	return str
}
