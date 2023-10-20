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
	dumpCmd.Flags().StringVarP(&dpFields, "fields", "f", "", "Specify the fields to be dumped from the input file. Default is all fields.")
	dumpCmd.Flags().StringVarP(&dpRewriterFiles, "rewrite-files", "r", "", "List of files that need to be rewritten based on the given configurations.")

	// input & output
	dumpCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "Path to the IP database file.")
	dumpCmd.Flags().StringVarP(&inputFormat, "input-format", "", "", "Specify the format of the input file. Examples: ipdb, mmdb, etc.")
	dumpCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Path to the dumped file.")
}

var dumpCmd = &cobra.Command{
	Use:   "dump -i inputFile [--input-format] [-o outputFile]",
	Short: "Dump data from IP database file to plain file.",
	Long: `Dump data from IP database file to plain file.

Example:
    ips dump -i geoip.mmdb -o geoip.txt
`,
	PreRun: PreRunInit,
	Run:    Dump,
}

func Dump(cmd *cobra.Command, args []string) {

	if len(args) == 0 && len(inputFile) == 0 {
		_ = cmd.Help()
		return
	}

	if len(inputFile) == 0 {
		inputFile = args[0]
	}

	if err := manager.Pack(inputFormat, inputFile, plain.DBFormat, outputFile); err != nil {
		log.Fatal(err)
	}
}
