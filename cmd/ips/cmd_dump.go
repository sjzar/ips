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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sjzar/ips/format/plain"
)

func init() {
	rootCmd.AddCommand(dumpCmd)

	// operate
	dumpCmd.Flags().StringVarP(&dpFields, "fields", "f", "", UsageDPFields)
	dumpCmd.Flags().StringVarP(&dpRewriterFiles, "rewrite-files", "r", "", UsageRewriteFiles)
	dumpCmd.Flags().StringVarP(&lang, "lang", "", "", UsageLang)

	// input & output
	dumpCmd.Flags().StringSliceVarP(&inputFile, "input-file", "i", nil, UsageDPInputFile)
	dumpCmd.Flags().StringSliceVarP(&inputFormat, "input-format", "", nil, UsageDPInputFormat)
	dumpCmd.Flags().StringVarP(&readerOption, "input-option", "", "", UsageReaderOption)
	dumpCmd.Flags().StringVarP(&hybridMode, "hybrid-mode", "", "aggregation", UsageHybridMode)
	dumpCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", UsageDumpOutputFile)
	dumpCmd.Flags().IntVarP(&readerJobs, "reader-jobs", "", 0, UsageReaderJobs)

}

var dumpCmd = &cobra.Command{
	Use:   "dump -i inputFile [--input-format] [-o outputFile]",
	Short: "Export IP database contents to a text file",
	Long: `Use the 'ips dump' command to extract and export data from IP databases into a plain text format, which can be tailored by specifying fields, formats, and languages.

For more detailed information and advanced configuration options, please refer to https://github.com/sjzar/ips/blob/main/docs/dump.md
`,
	Example: `  # Export all fields from an IP database file to a text file
  ips dump -i geoip.mmdb -o geoip.txt

  # Export specific fields (country and city) from an IP database file
  ips dump -i geoip.mmdb -o geoip.txt --fields "country,city"`,
	PreRun: PreRunInit,
	Run:    Dump,
}

func Dump(cmd *cobra.Command, args []string) {

	if len(args) == 0 && len(inputFile) == 0 {
		_ = cmd.Help()
		return
	}

	if len(inputFile) == 0 {
		inputFile = []string{args[0]}
	}

	if err := manager.Pack(inputFormat, inputFile, plain.DBFormat, outputFile); err != nil {
		log.Fatal(err)
	}
}
