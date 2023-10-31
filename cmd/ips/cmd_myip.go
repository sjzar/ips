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

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(myipCmd)

	// myip
	myipCmd.Flags().StringVarP(&localAddr, "local-addr", "", "", "Specify local address for outbound connections.")
	myipCmd.Flags().IntVarP(&myIPCount, "count", "", 0, "Number of detectors to confirm the IP.")
	myipCmd.Flags().IntVarP(&myIPTimeoutS, "timeout", "", 0, "Set the maximum wait time for detectors.")

	// operate
	myipCmd.Flags().StringVarP(&fields, "fields", "f", "", "Specify the fields of interest for the IP data. Separate multiple fields with commas.")
	myipCmd.Flags().BoolVarP(&useDBFields, "use-db-fields", "", false, "Use field names as they appear in the database. Default is common field names.")
	myipCmd.Flags().StringVarP(&rewriteFiles, "rewrite-files", "r", "", "List of files that need to be rewritten based on the given configurations.")
	myipCmd.Flags().StringVarP(&lang, "lang", "", "", "Set the language for the output. Example values: en, zh-CN, etc.")

	// database
	myipCmd.Flags().StringVarP(&rootFile, "file", "i", "", "Path to the IPv4 and IPv6 database file.")
	myipCmd.Flags().StringVarP(&rootFormat, "format", "", "", "Specify the format of the database. Examples: ipdb, mmdb, etc.")
	myipCmd.Flags().StringVarP(&rootIPv4File, "ipv4-file", "", "", "Path to the IPv4 database file.")
	myipCmd.Flags().StringVarP(&rootIPv4Format, "ipv4-format", "", "", "Specify the format for IPv4 data. Examples: ipdb, mmdb, etc.")
	myipCmd.Flags().StringVarP(&rootIPv6File, "ipv6-file", "", "", "Path to the IPv6 database file.")
	myipCmd.Flags().StringVarP(&rootIPv6Format, "ipv6-format", "", "", "Specify the format for IPv6 data. Examples: ipdb, mmdb, etc.")
	myipCmd.Flags().StringVarP(&readerOption, "database-option", "", "", "Specify the option for database reader.")

	// output
	myipCmd.Flags().StringVarP(&rootTextFormat, "text-format", "", "", "Specify the desired format for text output. It supports %origin and %values parameters.")
	myipCmd.Flags().StringVarP(&rootTextValuesSep, "text-values-sep", "", "", "Specify the separator for values in text output. (default is space)")
	myipCmd.Flags().BoolVarP(&rootJson, "json", "j", false, "Output the results in JSON format.")
	myipCmd.Flags().BoolVarP(&rootJsonIndent, "json-indent", "", false, "Output the results in indent JSON format.")
	myipCmd.Flags().BoolVarP(&rootAlfred, "alfred", "", false, "Output the results in Alfred format.")

}

var myipCmd = &cobra.Command{
	Use:   "myip",
	Short: "Retrieve your public IP address.",
	Long: `The 'myip' command uses multiple detectors to discover and return your public IP address. 
It is designed to work in scenarios with multiple network interfaces and allows you to specify 
the local address for outbound connections, the number of detectors to confirm the IP, and the timeout 
for the detection process.`,
	PreRun: PreRunInit,
	Run:    MyIP,
}

func MyIP(cmd *cobra.Command, args []string) {
	ip, err := manager.MyIP()
	if err != nil {
		return
	}
	fmt.Println(ip)
}
