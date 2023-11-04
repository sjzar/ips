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
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sjzar/ips/internal/ips"
)

var manager *ips.Manager

// Initialization function for setting up command line arguments and flags.
func init() {
	// common
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "", "info", UsageLogLevel)

	// operate

	rootCmd.Flags().StringVarP(&fields, "fields", "f", "", UsageFields)
	rootCmd.Flags().BoolVarP(&useDBFields, "use-db-fields", "", false, UsageUseDBFields)
	rootCmd.Flags().StringVarP(&rewriteFiles, "rewrite-files", "r", "", UsageRewriteFiles)
	rootCmd.Flags().StringVarP(&lang, "lang", "", "", UsageLang)

	// database
	rootCmd.Flags().StringVarP(&rootFile, "file", "i", "", UsageQueryFile)
	rootCmd.Flags().StringVarP(&rootFormat, "format", "", "", UsageQueryFormat)
	rootCmd.Flags().StringVarP(&rootIPv4File, "ipv4-file", "", "", UsageQueryIPv4File)
	rootCmd.Flags().StringVarP(&rootIPv4Format, "ipv4-format", "", "", UsageQueryIPv4Format)
	rootCmd.Flags().StringVarP(&rootIPv6File, "ipv6-file", "", "", UsageQueryIPv6File)
	rootCmd.Flags().StringVarP(&rootIPv6Format, "ipv6-format", "", "", UsageQueryIPv6Format)
	rootCmd.Flags().StringVarP(&readerOption, "database-option", "", "", UsageReaderOption)

	// output
	rootCmd.Flags().StringVarP(&rootTextFormat, "text-format", "", "", UsageTextFormat)
	rootCmd.Flags().StringVarP(&rootTextValuesSep, "text-values-sep", "", "", UsageTextValuesSep)
	rootCmd.Flags().BoolVarP(&rootJson, "json", "j", false, UsageJson)
	rootCmd.Flags().BoolVarP(&rootJsonIndent, "json-indent", "", false, UsageJsonIndent)
	rootCmd.Flags().BoolVarP(&rootAlfred, "alfred", "", false, UsageAlfred)
}

var rootCmd = &cobra.Command{
	Use:   "ips <ip or text>",
	Short: "IP Geolocation Database Tools",
	Long: `IP Geolocation Database Tools

The 'ips' is a command line tool for querying IP geolocation information and repacking database file.

It allows for flexible queries via command-line or pipe input, supporting both IPv4 and IPv6 formats, and provides customizable outputs.

For more detailed information and advanced configuration options, please refer to https://github.com/sjzar/ips/blob/main/docs/query.md
`,
	Example: `  # Standard IP query
  ips 8.8.8.8

  # Custom fields and output format
  ips 8.8.8.8 -f "country,city" --text-format "%values" --text-values-sep ":"

  # Pipeline query
  echo 8.8.8.8 | ips`,
	Args: cobra.MinimumNArgs(0),
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	PreRun: PreRunInit,
	Run:    Root,
}

// Execute is the entry point for the CLI tool. It executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}

// Root is the main logic for the IP query command. It also supports pipeline queries.
func Root(cmd *cobra.Command, args []string) {

	// Check for pipeline mode
	if len(args) == 0 {
		if fi, err := os.Stdin.Stat(); err != nil || fi.Mode()&os.ModeNamedPipe == 0 {
			_ = cmd.Help()
			return
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			if len(text) == 0 {
				continue
			}
			ret, err := manager.ParseText(text)
			if err != nil {
				log.Fatal(err)
			}
			if len(ret) == 0 {
				continue
			}
			fmt.Println(ret)
		}
		return
	}

	ret, err := manager.ParseText(strings.Join(args, " "))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ret)
}

// PreRunInit is called before the main Run function. It sets up logging and initializes the IP manager.
func PreRunInit(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			_, filename := path.Split(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
	lv, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(lv)
	if lv == log.DebugLevel {
		log.SetReportCaller(true)
	}

	// Initialize the IP manager with the config
	conf := GetFlagConfig()
	manager = ips.NewManager(conf)
}
