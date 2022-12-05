package main

import (
	"log"

	"github.com/sjzar/ips/cmd/ips"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ips.Execute()
}
