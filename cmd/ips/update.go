package ips

// IPDB
// Updated: 2018-11-18 (0.o)
// https://raw.githubusercontent.com/ipipdotnet/ipdb-go/master/city.free.ipdb

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
