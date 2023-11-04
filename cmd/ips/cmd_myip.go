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
	myipCmd.Flags().StringVarP(&fields, "fields", "f", "", UsageFields)
	myipCmd.Flags().BoolVarP(&useDBFields, "use-db-fields", "", false, UsageUseDBFields)
	myipCmd.Flags().StringVarP(&rewriteFiles, "rewrite-files", "r", "", UsageRewriteFiles)
	myipCmd.Flags().StringVarP(&lang, "lang", "", "", UsageLang)

	// database
	myipCmd.Flags().StringVarP(&rootFile, "file", "i", "", UsageQueryFile)
	myipCmd.Flags().StringVarP(&rootFormat, "format", "", "", UsageQueryFormat)
	myipCmd.Flags().StringVarP(&rootIPv4File, "ipv4-file", "", "", UsageQueryIPv4File)
	myipCmd.Flags().StringVarP(&rootIPv4Format, "ipv4-format", "", "", UsageQueryIPv4Format)
	myipCmd.Flags().StringVarP(&rootIPv6File, "ipv6-file", "", "", UsageQueryIPv6File)
	myipCmd.Flags().StringVarP(&rootIPv6Format, "ipv6-format", "", "", UsageQueryIPv6Format)
	myipCmd.Flags().StringVarP(&readerOption, "database-option", "", "", UsageReaderOption)

	// output
	myipCmd.Flags().StringVarP(&rootTextFormat, "text-format", "", "", UsageTextFormat)
	myipCmd.Flags().StringVarP(&rootTextValuesSep, "text-values-sep", "", "", UsageTextValuesSep)
	myipCmd.Flags().BoolVarP(&rootJson, "json", "j", false, UsageJson)
	myipCmd.Flags().BoolVarP(&rootJsonIndent, "json-indent", "", false, UsageJsonIndent)
	myipCmd.Flags().BoolVarP(&rootAlfred, "alfred", "", false, UsageAlfred)

}

var myipCmd = &cobra.Command{
	Use:   "myip",
	Short: "Retrieve your public IP address",
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
