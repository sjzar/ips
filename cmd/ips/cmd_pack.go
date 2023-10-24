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
	packCmd.Flags().StringVarP(&dpFields, "fields", "f", "", "Specify the fields to be dumped from the input file. Default is all fields.")
	packCmd.Flags().StringVarP(&dpRewriterFiles, "rewrite-files", "r", "", "List of files that need to be rewritten based on the given configurations.")
	packCmd.Flags().StringVarP(&lang, "lang", "", "", "Set the language for the output. Example values: en, zh-CN, etc.")

	// input & output
	packCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "Path to the IP database file.")
	packCmd.Flags().StringVarP(&inputFormat, "input-format", "", "", "Specify the format of the input file. Examples: ipdb, mmdb, etc.")
	packCmd.Flags().StringVarP(&readerOption, "input-option", "", "", "Specify the option for database reader.")
	packCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Path to the packed IP database file.")
	packCmd.Flags().StringVarP(&outputFormat, "output-format", "", "", "Specify the format of the output file. Examples: ipdb, mmdb, etc.")
	packCmd.Flags().StringVarP(&writerOption, "output-option", "", "", "Specify the option for database writer.")

}

var packCmd = &cobra.Command{
	Use:   "pack -i inputFile [--input-format format] -o outputFile [--output-format format]",
	Short: "Pack data from IP database file to another IP database file.",
	Long: `Pack data from IP database file to another IP database file.

Example:
    ips pack -i geoip.mmdb -o geoip.ipdb
`,
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
