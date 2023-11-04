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
	"log"
	"sort"

	"github.com/spf13/cobra"

	"github.com/sjzar/ips/internal/ips"
)

func init() {
	rootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "download [database_name] [custom_url]",
	Short: "Download IP database files",
	Long: `The 'ips download' command facilitates the acquisition and updating of IP geolocation database files.

For more detailed information and advanced configuration options, please refer to https://github.com/sjzar/ips/blob/main/docs/download.md
`,
	Example: `  # To download a predefined database file
  ips download city.free.ipdb

  # To download a database file from a custom URL
  ips download city.ipdb https://foo.com/city.ipdb

  # To configure the downloaded file as the default IPv4 database
  ips config set ipv4 city.ipdb`,
	PreRun: PreRunInit,
	Run:    Download,
}

func Download(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		fmt.Printf("IPS Download List =============================================================\n")
		list := make([]string, 0, len(ips.DownloadMap))
		for k := range ips.DownloadMap {
			list = append(list, k)
		}
		sort.Strings(list)
		for _, k := range list {
			fmt.Printf("%s: [%s]\n", k, ips.DownloadMap[k])
		}
		fmt.Println("===============================================================================")
		_ = cmd.Help()
		return
	}

	file, _url := "", ""
	if len(args) > 1 {
		file, _url = args[0], args[1]
	} else {
		file = args[0]
	}

	if err := manager.Download(file, _url); err != nil {
		log.Printf("Download %s failed: %s", args[0], err)
	}
}
