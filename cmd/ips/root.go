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
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sjzar/ips/parser"
)

var (
	rootDBFormat   string
	rootDBFile     string
	rootFields     string
	rootJsonFormat bool
	rootJsonIndent bool
)

func init() {
	rootCmd.Flags().StringVarP(&rootDBFormat, "format", "", "", "database format")
	rootCmd.Flags().StringVarP(&rootDBFile, "database", "d", "", "database file")
	rootCmd.Flags().StringVarP(&rootFields, "fields", "f", "", "fields")
	rootCmd.Flags().BoolVarP(&rootJsonFormat, "json", "", false, "json format")
	rootCmd.Flags().BoolVarP(&rootJsonIndent, "json-indent", "", false, "json indent")
}

var rootCmd = &cobra.Command{
	Use:   "ips ip [-d dbFile] [-f fields] [--json]",
	Short: "ips commandline tools",
	Long:  `ips is a tool for querying, scanning, and packing IP geolocation databases.`,
	Args:  cobra.MinimumNArgs(0),
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	Run: Root,
}

// Execute 用于启动命令行工具
func Execute() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Root IP查询命令, 支持pipeline查询
func Root(cmd *cobra.Command, args []string) {

	// pipeline mode
	if len(args) == 0 {
		fi, err := os.Stdin.Stat()
		if err != nil || fi.Mode()&os.ModeNamedPipe == 0 {
			_ = cmd.Help()
			return
		}
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			str := scanner.Text()
			if len(str) == 0 {
				continue
			}
			if rootJsonFormat {
				result := ParseLineJson(str, rootJsonIndent)
				if !strings.Contains(result, "[]") {
					fmt.Println(result)
				}
			} else {
				fmt.Println(ParseLine(str))
			}
		}
		return
	}

	if rootJsonFormat {
		fmt.Println(ParseLineJson(strings.Join(args, " "), rootJsonIndent))
	} else {
		fmt.Println(ParseLine(strings.Join(args, " ")))
	}
}

// ParseLine 解析文本
func ParseLine(line string) string {
	p := parser.NewTextParser(line)
	p.IPv4FillResult = func(str string) []string {
		_, values, err := GetIPv4().Find(net.ParseIP(str))
		if err != nil {
			return []string{}
		}
		return values
	}
	p.IPv6FillResult = func(str string) []string {
		_, values, err := GetIPv6().Find(net.ParseIP(str))
		if err != nil {
			return []string{}
		}
		return values
	}

	return p.Parse().String()
}

// ParseLineJson 解析文本并返回json格式
func ParseLineJson(line string, indent bool) string {
	p := parser.NewTextParser(line)
	p.IPv4Fields = GetIPv4().Meta().Fields
	p.IPv4FillResult = func(str string) []string {
		_, values, err := GetIPv4().Find(net.ParseIP(str))
		if err != nil {
			return []string{}
		}
		return values
	}
	p.IPv6Fields = GetIPv6().Meta().Fields
	p.IPv6FillResult = func(str string) []string {
		_, values, err := GetIPv6().Find(net.ParseIP(str))
		if err != nil {
			return []string{}
		}
		return values
	}
	return p.Parse().Json(indent)
}
