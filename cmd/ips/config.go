/*
 * Copyright (c) 2022 shenjunzheng@gmail.com
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

	"github.com/spf13/cobra"

	"github.com/sjzar/ips/cmd/ips/conf"
)

var FormatFileMap = map[string]string{
	"ipdb":      "city.free.ipdb",
	"qqwry":     "qqwry.dat",
	"zxinc":     "zxipv6wry.db",
	"geoip2":    "GeoLite2-City.mmdb",
	"ip2region": "ip2region.db",
	"dbip":      "dbip-city-lite.mmdb",
}

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config [set <key> <value>]",
	Short: "Print the config of ips",
	Long: `1. choose free ip database format
  free ipv4 format support: ipdb, qqwry, geoip2, ip2region, dbip
  free ipv6 format support: zxinc, geoip2, dbip
  use 'ips config set ipv4 ipdb' or 'ips config set ipv6 zxinc' to set format

2. set ip database file path
  use 'ips config set ipv4_file ~/path/to/ipv4.db' to set file path
  use 'ips config set ipv4_format format' to set database format

3. set print fields
  use 'ips config set fields country,province,city,isp' to set print fields`,
	Run: Config,
}

func Config(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("config_path: ", conf.ConfigPath)
		for k, v := range conf.GetConfig() {
			fmt.Printf("%s: %v\n", k, v)
		}
		fmt.Println("\n================================================================================")
		fmt.Println(cmd.Long)
		return
	}
	switch args[0] {
	case "set":
		if len(args) < 3 {
			log.Fatal("need key and value")
			return
		}
		key, value := args[1], args[2]

		switch key {
		case "ipv4":
			if dbfile, ok := FormatFileMap[value]; ok {
				key = "ipv4_file"
				value = dbfile
			}
		case "ipv6":
			if dbfile, ok := FormatFileMap[value]; ok {
				key = "ipv6_file"
				value = dbfile
			}
		}

		if err := conf.SetConfig(key, value); err != nil {
			log.Fatal(err)
		}
		fmt.Println("set config success")
	default:
		log.Fatal("unknown config command")
	}
}
