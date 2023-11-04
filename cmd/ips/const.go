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

const (
	// Flag Usage
	// Global Flags

	UsageLogLevel = "Set the desired verbosity level for logging."

	// Operate Flags

	UsageLang         = "Language for output data. (default \"zh-CN\")"
	UsageFields       = "Fields to include in the output, separated by commas. (default \"country,province,city,isp\")"
	UsageUseDBFields  = "Use field names as they appear in the database. Default is common field names."
	UsageRewriteFiles = "Paths to files containing data rewrite rules, separated by commas."
	UsageDPFields     = "Fields to extract from the database. Defaults to all available fields."

	// Database Flags

	UsageQueryFile        = "Path to the combined IPv4/IPv6 database file."
	UsageQueryFormat      = "The format of the IPv4/IPv6 database file."
	UsageQueryIPv4File    = "Path to the IPv4 database file."
	UsageQueryIPv4Format  = "The format of the IPv4 database file."
	UsageQueryIPv6File    = "Path to the IPv6 database file."
	UsageQueryIPv6Format  = "The format of the IPv6 database file."
	UsageDPInputFile      = "Path to the input IP database file (required)."
	UsageDPInputFormat    = "The format of the input IP database file."
	UsageDumpOutputFile   = "Destination path for the dumped data. Defaults to standard output if not specified."
	UsagePackOutputFile   = "Path to the output IP database file (required)."
	UsagePackOutputFormat = "The format for the output IP database file."
	UsageReaderOption     = "Additional options for the database reader, if applicable."
	UsageWriterOption     = "Additional options for the database writer, if applicable."

	// Output Flags

	UsageTextFormat    = "Specify the desired format for text output. (default \"%origin [%values]\")"
	UsageTextValuesSep = "Specify the separator for values in text output. (default \" \")"
	UsageJson          = "Output the results in JSON format."
	UsageJsonIndent    = "Output the results in indent JSON format."
	UsageAlfred        = "Output the results in Alfred format."
)
