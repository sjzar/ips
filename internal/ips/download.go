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
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/sjzar/ips/pkg/errors"

	"github.com/sjzar/ips/internal/util"
)

/*
Download Database Files

IPDB
Official:
City Free (2018-11-18): https://raw.githubusercontent.com/ipipdotnet/ipdb-go/master/city.free.ipdb

CZ88.NET
Official:
qqwry.dat (daily update): https://github.com/HMBSbige/qqwry/releases/download/2023-10-18/qqwry.dat
qqwry.dat Mirror 2 (2022-04-20): https://github.com/metowolf/qqwry.dat/releases/download/20231011/qqwry.dat

MaxMind
Official:
GeoLite2-City.mmdb (2023-10-16): https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-City.mmdb (https://git.io/GeoLite2-City.mmdb)

ZX Inc.
Official:
zxipv6wry.db (2021-05-11): https://raw.githubusercontent.com/ZX-Inc/zxipdb-python/main/data/ipv6wry.db
zxipv6wry.db Mirror 2 (2021-05-11): https://ip.zxinc.org/ip.7z

ip2region
Official:
ip2region.db (2022-12-07): https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb

DB-IP
Official:
dbip-city-lite.mmdb (2023-10): https://download.db-ip.com/free/dbip-city-lite-2023-10.mmdb.gz
dbip-asn-lite.mmdb (2023-10): https://download.db-ip.com/free/dbip-asn-lite-2023-10.mmdb.gz

*/

var DownloadMap = map[string]string{
	"city.free.ipdb":      "https://raw.githubusercontent.com/ipipdotnet/ipdb-go/master/city.free.ipdb",
	"qqwry.dat":           "https://github.com/metowolf/qqwry.dat/releases/download/20231011/qqwry.dat",
	"zxipv6wry.db":        "https://raw.githubusercontent.com/ZX-Inc/zxipdb-python/main/data/ipv6wry.db",
	"GeoLite2-City.mmdb":  "https://git.io/GeoLite2-City.mmdb",
	"ip2region.xdb":       "https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb",
	"dbip-city-lite.mmdb": "https://download.db-ip.com/free/dbip-city-lite-2023-10.mmdb.gz",
	"dbip-asn-lite.mmdb":  "https://download.db-ip.com/free/dbip-asn-lite-2023-10.mmdb.gz",
}

// Download downloads the database file to ips dir.
func (m *Manager) Download(file, _url string) error {

	if len(_url) == 0 {
		_url = DownloadMap[file]
		if len(_url) == 0 {
			log.Debugf("unknown file %s", file)
			return errors.ErrFileNotFound
		}
	}

	fmt.Printf("Downloading %s from %s to %s\n", file, _url, m.Conf.IPSDir)

	f, err := os.Create(m.Conf.IPSDir + "/" + file)
	if err != nil {
		log.Debugf("create file %s failed: %s", file, err)
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	resp, err := http.DefaultClient.Get(_url)
	if err != nil {
		log.Debugf("http get %s failed: %s", _url, err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		log.Debugf("http get %s failed: %s", _url, resp.Status)
		return errors.ErrFailedDownload
	}

	var r io.Reader = resp.Body
	if u, _ := url.Parse(_url); filepath.Ext(u.Path) == ".gz" {
		if r, err = gzip.NewReader(resp.Body); err != nil {
			log.Debugf("gzip.NewReader failed: %s", err)
			return err
		}
	}

	bar := util.ProgressBar(
		resp.ContentLength,
		"downloading",
	)
	_, err = io.Copy(io.MultiWriter(f, bar), r)
	if err != nil {
		log.Debugf("io.Copy failed: %s", err)
		return err
	}

	fmt.Println("Download " + file + " success.")
	return nil
}
