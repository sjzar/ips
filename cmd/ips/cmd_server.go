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
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	// server
	serverCmd.Flags().StringVarP(&addr, "addr", "a", "", "server listen address")

	// operate
	serverCmd.Flags().StringVarP(&fields, "fields", "f", "", "Specify the fields of interest for the IP data. Separate multiple fields with commas.")
	serverCmd.Flags().BoolVarP(&useDBFields, "use-db-fields", "", false, "Use field names as they appear in the database. Default is common field names.")
	serverCmd.Flags().StringVarP(&rewriteFiles, "rewrite-files", "r", "", "List of files that need to be rewritten based on the given configurations.")
	serverCmd.Flags().StringVarP(&lang, "lang", "", "", "Set the language for the output. Example values: en, zh-CN, etc.")

	// database
	serverCmd.Flags().StringVarP(&rootFile, "file", "i", "", "Path to the IPv4 and IPv6 database file.")
	serverCmd.Flags().StringVarP(&rootFormat, "format", "", "", "Specify the format of the database. Examples: ipdb, mmdb, etc.")
	serverCmd.Flags().StringVarP(&rootIPv4File, "ipv4-file", "", "", "Path to the IPv4 database file.")
	serverCmd.Flags().StringVarP(&rootIPv4Format, "ipv4-format", "", "", "Specify the format for IPv4 data. Examples: ipdb, mmdb, etc.")
	serverCmd.Flags().StringVarP(&rootIPv6File, "ipv6-file", "", "", "Path to the IPv6 database file.")
	serverCmd.Flags().StringVarP(&rootIPv6Format, "ipv6-format", "", "", "Specify the format for IPv6 data. Examples: ipdb, mmdb, etc.")
	serverCmd.Flags().StringVarP(&readerOption, "database-option", "", "", "Specify the option for database reader.")
}

var serverCmd = &cobra.Command{
	Use:    "server",
	Short:  "Start IPS server",
	PreRun: PreRunInit,
	Run:    Server,
}

func Server(cmd *cobra.Command, args []string) {
	manager.Service()
}
