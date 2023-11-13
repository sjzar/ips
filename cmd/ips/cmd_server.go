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
	serverCmd.Flags().StringVarP(&addr, "addr", "a", "", "Listen address")

	// operate
	serverCmd.Flags().StringVarP(&fields, "fields", "f", "", UsageFields)
	serverCmd.Flags().BoolVarP(&useDBFields, "use-db-fields", "", false, UsageUseDBFields)
	serverCmd.Flags().StringVarP(&rewriteFiles, "rewrite-files", "r", "", UsageRewriteFiles)
	serverCmd.Flags().StringVarP(&lang, "lang", "", "", UsageLang)

	// database
	serverCmd.Flags().StringSliceVarP(&rootFile, "file", "i", nil, UsageQueryFile)
	serverCmd.Flags().StringSliceVarP(&rootFormat, "format", "", nil, UsageQueryFormat)
	serverCmd.Flags().StringSliceVarP(&rootIPv4File, "ipv4-file", "", nil, UsageQueryIPv4File)
	serverCmd.Flags().StringSliceVarP(&rootIPv4Format, "ipv4-format", "", nil, UsageQueryIPv4Format)
	serverCmd.Flags().StringSliceVarP(&rootIPv6File, "ipv6-file", "", nil, UsageQueryIPv6File)
	serverCmd.Flags().StringSliceVarP(&rootIPv6Format, "ipv6-format", "", nil, UsageQueryIPv6Format)
	serverCmd.Flags().StringVarP(&readerOption, "database-option", "", "", UsageReaderOption)
	serverCmd.Flags().StringVarP(&hybridMode, "hybrid-mode", "", "aggregation", UsageHybridMode)

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
