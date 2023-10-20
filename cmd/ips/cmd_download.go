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
	Use:   "download [file] [url]",
	Short: "Download database files",
	Long: `Example:
 # download list file city.free.ipdb
 ips download city.free.ipdb

 # download another database file
 ips download city.ipdb https://foo.com/city.ipdb

 # set database file after download
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
