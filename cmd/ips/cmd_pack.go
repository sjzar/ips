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
)

func init() {
	rootCmd.AddCommand(packCmd)

	// operate
	packCmd.Flags().StringVarP(&dpFields, "fields", "f", "", UsageDPFields)
	packCmd.Flags().StringVarP(&dpRewriterFiles, "rewrite-files", "r", "", UsageRewriteFiles)
	packCmd.Flags().StringVarP(&lang, "lang", "", "", UsageLang)

	// input & output
	packCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", UsageDPInputFile)
	packCmd.Flags().StringVarP(&inputFormat, "input-format", "", "", UsageDPInputFormat)
	packCmd.Flags().StringVarP(&readerOption, "input-option", "", "", UsageReaderOption)
	packCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", UsagePackOutputFile)
	packCmd.Flags().StringVarP(&outputFormat, "output-format", "", "", UsagePackOutputFormat)
	packCmd.Flags().StringVarP(&writerOption, "output-option", "", "", UsageWriterOption)

}

var packCmd = &cobra.Command{
	Use:   "pack -i inputFile [--input-format format] -o outputFile [--output-format format]",
	Short: "Repackage IP database file",
	Long: `The 'ips pack' command enables users to create a new IP database file from an existing one while allowing for customization of the data fields included.

For more detailed information and advanced configuration options, please refer to https://github.com/sjzar/ips/blob/main/docs/pack.md
`,
	Example: `  # Package IP Database and Specify Fields
  ips pack -i geoip.mmdb -o geoip_custom.ipdb --fields "country,city"`,
	PreRun: PreRunInit,
	Run:    Pack,
}

func Pack(cmd *cobra.Command, args []string) {

	if len(inputFile) == 0 || len(outputFile) == 0 {
		_ = cmd.Help()
		return
	}

	if err := manager.Pack(inputFormat, inputFile, outputFormat, outputFile); err != nil {
		log.Fatal(err)
	}
}
