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
	Short: "View or modify ips configuration items",
	Long: `Example:
 # set ipv4 database file path
 ips config set ipv4 ~/path/to/ipv4.db

 # set ipv4 database format
 ips config set ipv4_format ipdb

 # set query fields
 ips config set fields "country,province,city,isp"

 # unset ipv6 database file path
 ips config unset ipv6

 # reset config
 ips config reset

 # ipv4 format: ipip, qqwry, maxmind, ip2region, dbip
 # ipv6 format: zxinc, maxmind, dbip
 # use 'ips config set ipv4 ipdb' or 'ips config set ipv6 zxinc' to set format`,
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
		SetConfig(args[1], args[2])
	case CmdUnset:
		if len(args) < 2 {
			log.Fatal("missing key")
			return
		}
		SetConfig(args[1], "")
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

func SetConfig(key string, value string) {

	if _key, ok := MagicMap[key]; ok {
		key = _key
	}

	if _value, ok := MagicMap[value]; ok {
		value = _value
	}

	if err := config.SetConfig(key, value); err != nil {
		log.Fatal(err)
	}

	if len(value) == 0 {
		log.Infof("unset %s success", key)
	} else {
		log.Infof("set %s: [%s] success", key, value)
	}
}
