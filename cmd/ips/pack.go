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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sjzar/ips/data"
	"github.com/sjzar/ips/db"
	"github.com/sjzar/ips/db/ipdb"
	"github.com/sjzar/ips/ipio"
	"github.com/sjzar/ips/model"
	"github.com/sjzar/ips/rewriter"
)

var (
	packDBFormat     string
	packFields       string
	packOutput       string
	packRewriteFiles string
	packIPv4         bool
	packIPv6         bool
)

func init() {
	rootCmd.AddCommand(packCmd)
	packCmd.Flags().StringVarP(&packDBFormat, "format", "", "", "database format")
	packCmd.Flags().StringVarP(&packFields, "fields", "f", "", "fields")
	packCmd.Flags().StringVarP(&packOutput, "output", "o", "", "output")
	packCmd.Flags().StringVarP(&packRewriteFiles, "rewrite", "r", "", "rewrite files")
	packCmd.Flags().BoolVarP(&packIPv4, "ipv4", "", false, "ipv4 only")
	packCmd.Flags().BoolVarP(&packIPv6, "ipv6", "", false, "ipv6 only")
}

var packCmd = &cobra.Command{
	Use:   "pack [-f fields] [-r rewrite] [-o output] [--ipv4 | --ipv6] [--format dbFormat] scanFile",
	Short: "Packing IP geolocation database file",
	Long:  `Packing IP geolocation database file`,
	Run:   Pack,
}

func Pack(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Usage()
		return
	}

	var s ipio.Scanner
	var err error
	ext := filepath.Ext(args[0])
	if len(packDBFormat) != 0 || ext == ".ipdb" || ext == ".awdb" {
		database, err := db.NewDatabase(packDBFormat, args[0])
		if err != nil {
			log.Fatal("new database failed ", err)
		}
		fields := strings.Join(database.Meta().Fields, ",")
		if len(packFields) != 0 {
			fields = packFields
		}
		selector := ipio.NewFieldSelector(fields)
		rw := rewriter.NewDataRewriter(nil, nil)
		rw.DataLoader.LoadString(data.ASN2ISP)
		rw.DataLoader.LoadString(data.Province, data.City, data.ISP)
		if len(packRewriteFiles) != 0 {
			rw = rewriter.NewDataRewriter(nil, rw)
			split := strings.Split(packRewriteFiles, ",")
			for i := range split {
				if err := rw.DataLoader.LoadFile(split[i]); err != nil {
					log.Println("load rewrite file failed", err)
				}
			}
		}
		s = ipio.NewDBScanner(database, selector, rw)
	} else {
		s, err = ipio.NewIPScanScanner(args[0])
		if err != nil {
			log.Fatal(err)
		}
	}

	defer func() {
		_ = s.Close()
	}()

	if len(packOutput) == 0 {
		path, file := filepath.Split(args[0])
		fileName := strings.TrimSuffix(file, filepath.Ext(file))
		packOutput = path + fileName + ".ipdb"
	}

	ipVersion := uint16(0)
	if packIPv4 {
		ipVersion |= model.IPv4
	}
	if packIPv6 {
		ipVersion |= model.IPv6
	}

	meta := s.Meta()
	if ipVersion != 0 {
		meta.IPVersion = ipVersion
	}

	w := ipdb.NewWriter(s.Meta(), nil)
	if err := ipio.ScanWrite(s, w); err != nil {
		log.Fatal(err)
	}

	if len(packOutput) == 0 {
		packOutput = strings.TrimSuffix(args[0], filepath.Ext(args[0])) + ".ipdb"
	}

	if packOutput == args[0] {
		packOutput = strings.TrimSuffix(packOutput, filepath.Ext(packOutput)) + ".pack.ipdb"
	}

	f, err := os.OpenFile(packOutput, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = f.Close()
	}()

	if err := w.Save(f); err != nil {
		log.Fatal(err)
	}
}
