package ips

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sjzar/ips/db"
	"github.com/sjzar/ips/ipio"

	"github.com/sjzar/ips/cmd/ips/conf"
)

var (
	// IPv4
	ipv4     ipio.Reader
	ipv4once sync.Once

	// IPv6
	ipv6     ipio.Reader
	ipv6once sync.Once
)

// GetIPv4 returns a ipio.Reader for IPv4
func GetIPv4() ipio.Reader {
	ipv4once.Do(func() {
		format := rootDBFormat
		if len(format) == 0 {
			format = conf.Conf.IPv4Format
		}
		file := rootDBFile
		if len(file) == 0 {
			file = conf.ConfigPath + "/" + conf.Conf.IPv4File
		}
		database, err := db.NewDatabase(format, file)
		if err != nil {
			log.Println("read database failed", file, err)
			if file == conf.ConfigPath+"/"+conf.Conf.IPv4File {
				log.Println("use ips update first")
			}
			os.Exit(1)
		}
		fields := conf.Conf.Fields
		if len(rootFields) != 0 {
			fields = strings.Split(rootFields, ",")
		}
		if len(fields) == 0 {
			fields = database.Meta().Fields
		}
		selector := ipio.NewFieldSelector(strings.Join(fields, ","))
		ipv4 = ipio.NewDBScanner(database, selector, nil)
	})

	return ipv4
}

// GetIPv6 returns a ipio.Reader for IPv6
func GetIPv6() ipio.Reader {
	ipv6once.Do(func() {
		format := rootDBFormat
		if len(format) == 0 {
			format = conf.Conf.IPv4Format
		}
		file := rootDBFile
		if len(file) == 0 {
			file = conf.ConfigPath + "/" + conf.Conf.IPv6File
		}
		database, err := db.NewDatabase(format, file)
		if err != nil {
			log.Println("read database failed", file, err)
			if file == conf.ConfigPath+"/"+conf.Conf.IPv6File {
				log.Println("use ips update first")
			}
			os.Exit(1)
		}
		fields := conf.Conf.Fields
		if len(rootFields) != 0 {
			fields = strings.Split(rootFields, ",")
		}
		if len(fields) == 0 {
			fields = database.Meta().Fields
		}
		selector := ipio.NewFieldSelector(strings.Join(fields, ","))
		ipv6 = ipio.NewDBScanner(database, selector, nil)
	})

	return ipv6
}
