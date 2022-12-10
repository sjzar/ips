package ips

// IPDB
// Updated: 2018-11-18 (0.o)
// https://raw.githubusercontent.com/ipipdotnet/ipdb-go/master/city.free.ipdb

// CZ88.NET
// Updated: 2022-04-20
// https://raw.githubusercontent.com/out0fmemory/qqwry.dat/master/historys/2022_04_20/qqwry.dat

// CZ88.NET Mirror 2
// Updated: daily update
// https://99wry.cf/qqwry.dat

// MaxMind GeoLite2-City
// Updated: 2022-12-07
// https://git.io/GeoLite2-City.mmdb
// https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-City.mmdb

// ZXINC
// Updated: 2021-05-11
// https://ip.zxinc.org/ip.7z
// https://raw.githubusercontent.com/ZX-Inc/zxipdb-python/main/data/ipv6wry.db

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/sjzar/ips/cmd/ips/conf"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update database",
	Run:   Update,
}

func Update(cmd *cobra.Command, args []string) {
	Download("city.free.ipdb", "https://raw.githubusercontent.com/ipipdotnet/ipdb-go/master/city.free.ipdb")
	Download("qqwry.dat", "https://99wry.cf/qqwry.dat")
	Download("zxipv6wry.db", "https://raw.githubusercontent.com/ZX-Inc/zxipdb-python/main/data/ipv6wry.db")

	//Download("GeoLite2-City.mmdb", "https://git.io/GeoLite2-City.mmdb")
}

func Download(file, url string) {
	log.Println("Downloading " + file + "...")
	f, err := os.OpenFile(conf.ConfigPath+"/"+file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("open file failed ", err)
	}
	defer f.Close()

	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Fatal("download file failed ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("download file failed ", resp.Status)
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatal("download file failed ", err)
	}
	log.Println("Download " + file + " success.")
}
