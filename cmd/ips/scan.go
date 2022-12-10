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
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sjzar/ips/data"
	"github.com/sjzar/ips/db"
	"github.com/sjzar/ips/ipio"
	"github.com/sjzar/ips/model"
	"github.com/sjzar/ips/rewriter"
)

var (
	scanDBFormat     string
	scanOutput       string
	scanFields       string
	scanRewriteFiles string
	scanIPv4         bool
	scanIPv6         bool
)

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&scanDBFormat, "format", "", "", "database format")
	scanCmd.Flags().StringVarP(&scanFields, "fields", "f", "", "fields")
	scanCmd.Flags().StringVarP(&scanOutput, "output", "o", "", "write to file instead of stdout")
	scanCmd.Flags().StringVarP(&scanRewriteFiles, "rewrite", "r", "", "rewrite files")
	scanCmd.Flags().BoolVarP(&scanIPv4, "ipv4", "", false, "ipv4 only")
	scanCmd.Flags().BoolVarP(&scanIPv6, "ipv6", "", false, "ipv6 only")
}

var scanCmd = &cobra.Command{
	Use:   "scan [-f fields] [-r rewrite] [-o output] [--ipv4 | --ipv6] [--format dbFormat] dbFile",
	Short: "Scanning IP geolocation database file",
	Long: `Scanning IP geolocation database file
`,
	Run: Scan,
}

func Scan(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Usage()
		return
	}

	ipVersion := uint16(0)
	if scanIPv4 {
		ipVersion |= model.IPv4
	}
	if scanIPv6 {
		ipVersion |= model.IPv6
	}

	database, err := db.NewDatabase(scanDBFormat, args[0])
	if err != nil {
		log.Fatal("new database failed ", err)
	}
	fields := strings.Join(database.Meta().Fields, ",")
	if len(scanFields) != 0 {
		fields = scanFields
	}
	selector := ipio.NewFieldSelector(fields)
	rw := rewriter.NewDataRewriter(nil, nil)
	rw.DataLoader.LoadString(data.ASN2ISP)
	rw.DataLoader.LoadString(data.Province, data.City, data.ISP)
	if len(scanRewriteFiles) != 0 {
		rw = rewriter.NewDataRewriter(nil, rw)
		split := strings.Split(scanRewriteFiles, ",")
		for i := range split {
			if err := rw.DataLoader.LoadFile(split[i]); err != nil {
				log.Fatal("load rewrite file failed ", err)
			}
		}
	}
	r := ipio.NewDBScanner(database, selector, rw)

	meta := r.Meta()
	if ipVersion != 0 {
		meta.IPVersion = ipVersion
	}
	if err := r.Init(meta); err != nil {
		log.Fatal("init scanner failed ", err)
	}
	output := os.Stdout
	if len(scanOutput) != 0 {
		output, err = os.Create(scanOutput)
		if err != nil {
			log.Fatal("create file failed ", err)
		}
	}

	w, err := ipio.NewIPScanWriter(r.Meta(), output)
	if err != nil {
		log.Fatal("new writer failed ", err)
	}
	if err := ipio.ScanWrite(r, w); err != nil {
		log.Fatal("scan write failed ", err)
	}
}
