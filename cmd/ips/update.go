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
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sjzar/ips/cmd/ips/conf"
)

/*
Download Database Files

IPDB
Official:
City Free (2018-11-18): https://raw.githubusercontent.com/ipipdotnet/ipdb-go/master/city.free.ipdb

CZ88.NET
Official:
qqwry.dat (daily update): https://99wry.cf/qqwry.dat
qqwry.dat Mirror 2 (2022-04-20): https://raw.githubusercontent.com/out0fmemory/qqwry.dat/master/historys/2022_04_20/qqwry.dat

MaxMind
Official:
GeoLite2-City.mmdb (2022-12-13): https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-City.mmdb (https://git.io/GeoLite2-City.mmdb)

ZX Inc.
Official:
zxipv6wry.db (2021-05-11): https://raw.githubusercontent.com/ZX-Inc/zxipdb-python/main/data/ipv6wry.db
zxipv6wry.db Mirror 2 (2021-05-11): https://ip.zxinc.org/ip.7z

ip2region
Official:
ip2region.db (2022-12-07): https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb

DB-IP
Official:
dbip-city-lite.mmdb (2022-12): https://download.db-ip.com/free/dbip-city-lite-2022-12.mmdb.gz
dbip-asn-lite.mmdb (2022-12): https://download.db-ip.com/free/dbip-asn-lite-2022-12.mmdb.gz

*/

var DownloadMap = map[string]string{
	"city.free.ipdb":      "https://raw.githubusercontent.com/ipipdotnet/ipdb-go/master/city.free.ipdb",
	"qqwry.dat":           "https://99wry.cf/qqwry.dat",
	"zxipv6wry.db":        "https://raw.githubusercontent.com/ZX-Inc/zxipdb-python/main/data/ipv6wry.db",
	"GeoLite2-City.mmdb":  "https://git.io/GeoLite2-City.mmdb",
	"ip2region.xdb":       "https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb",
	"dbip-city-lite.mmdb": "https://download.db-ip.com/free/dbip-city-lite-2022-12.mmdb.gz",
	"dbip-asn-lite.mmdb":  "https://download.db-ip.com/free/dbip-asn-lite-2022-12.mmdb.gz",
}

var (
	updateFile string
	updateURL  string
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&updateFile, "file", "f", "", "update file")
	updateCmd.Flags().StringVarP(&updateURL, "url", "u", "", "update url")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update database",
	Run:   Update,
}

func Update(cmd *cobra.Command, args []string) {
	if len(updateFile) > 0 {
		Download(updateFile, updateURL)
	} else {
		update()
	}
}

func update() {
	Download(conf.Conf.IPv4File, "")
	Download(conf.Conf.IPv6File, "")
}

func Download(file, _url string) {
	if len(_url) == 0 {
		_url = DownloadMap[file]
		if len(_url) == 0 {
			return
		}
	}

	log.Println("Downloading " + file + "...")
	f, err := os.OpenFile(conf.ConfigPath+"/"+file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("open file failed ", err)
	}
	defer f.Close()

	resp, err := http.DefaultClient.Get(_url)
	if err != nil {
		log.Fatal("download file failed ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("download file failed ", resp.Status)
	}

	var r io.Reader = resp.Body
	u, _ := url.Parse(_url)
	if filepath.Ext(u.Path) == ".gz" {
		if r, err = gzip.NewReader(resp.Body); err != nil {
			log.Fatal("gzip.NewReader failed ", err)
		}
	}

	_, err = io.Copy(f, r)
	if err != nil {
		log.Fatal("download file failed ", err)
	}
	log.Println("Download " + file + " success.")
}
