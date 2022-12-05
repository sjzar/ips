package ips

import (
	"log"
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
		database, err := db.NewDatabase(conf.Conf.IPv4Format, conf.ConfigPath+"/"+conf.Conf.IPv4File)
		if err != nil {
			log.Fatal(err)
		}
		fields := conf.Conf.Fields
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
		database, err := db.NewDatabase(conf.Conf.IPv6Format, conf.ConfigPath+"/"+conf.Conf.IPv6File)
		if err != nil {
			log.Fatal("read database failed ", conf.Conf.IPv4File, err)
		}
		fields := conf.Conf.Fields
		if len(fields) == 0 {
			fields = database.Meta().Fields
		}
		selector := ipio.NewFieldSelector(strings.Join(fields, ","))
		ipv6 = ipio.NewDBScanner(database, selector, nil)
	})

	return ipv6
}
