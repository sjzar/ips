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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sjzar/ips/internal/config"
)

const (
	CmdSet   = "set"
	CmdUnset = "unset"
	CmdReset = "reset"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config [set <key> <value>] [unset <key>] [reset]",
	Short: "Manage IPS configuration settings",
	Long: `The 'ips config' command allows you to view and change various settings for IPS.

You can specify database file paths, select output fields, and choose the format for IP databases.

Use 'set' to change a setting, 'unset' to remove it, and 'reset' to revert all settings to default values.

For IPv6, use similar commands, like 'ips config set ipv6_format mmdb' to specify the format.

For more detailed information and advanced configuration options, please refer to https://github.com/sjzar/ips/blob/main/docs/config.md
`,
	Example: `  # Set the IPv4 database file path:
  ips config set ipv4 ~/path/to/ipv4.db

  # Set the IPv4 database format:
  ips config set ipv4_format ipdb

  # Configure the fields to display in query results:
  ips config set fields "country,province,city,isp"

  # Remove the IPv6 database file path setting:
  ips config unset ipv6

  # Reset all configuration settings to their default values:
  ips config reset`,
	PreRun: PreRunInit,
	Run:    Config,
}

func Config(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		conf := GetConfig()
		fmt.Printf("IPS CONFIG ====================================================================\n\n")
		fmt.Println(conf.ShowConfig(false))
		fmt.Println("===============================================================================")
		_ = cmd.Help()
		return
	}
	switch args[0] {
	case CmdSet:
		if len(args) < 3 {
			log.Fatal("missing key or value")
			return
		}
		SetConfig(args[1], args[2:])
	case CmdUnset:
		if len(args) < 2 {
			log.Fatal("missing key")
			return
		}
		SetConfig(args[1], []string{})
	case CmdReset:
		if err := config.ResetConfig(); err != nil {
			log.Fatal("reset config failed ", err)
		}
		log.Info("reset config success")
	default:
		log.Fatal("unknown config command")
	}
}

var MagicMap = map[string]string{
	"ipv4":      "ipv4_file",
	"ipv6":      "ipv6_file",
	"ipip":      "city.free.ipdb",
	"qqwry":     "qqwry.dat",
	"zxinc":     "zxipv6wry.db",
	"maxmind":   "GeoLite2-City.mmdb",
	"ip2region": "ip2region.xdb",
	"dbip":      "dbip-city-lite.mmdb",
	"dbip-asn":  "dbip-asn-lite.mmdb",
}

func SetConfig(key string, value []string) {

	if _key, ok := MagicMap[key]; ok {
		key = _key
	}

	for i := range value {
		if _value, ok := MagicMap[value[i]]; ok {
			value[i] = _value
		}
	}

	var val interface{}
	val = value
	if len(value) == 1 {
		val = value[0]
	}

	if err := config.SetConfig(key, val); err != nil {
		log.Fatal(err)
	}

	if len(value) == 0 {
		log.Infof("unset %s success", key)
	} else {
		log.Infof("set %s: [%s] success", key, strings.Join(value, ","))
	}
}
