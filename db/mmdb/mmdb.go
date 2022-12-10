package mmdb

import (
	"net"

	"github.com/oschwald/maxminddb-golang"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

type Database struct {
	meta model.Meta
	db   *maxminddb.Reader
}

// New 初始化 mmdb 数据库实例
// 发现 GeoLite2-City 有数据不全的情况
func New(file string) (*Database, error) {

	db, err := maxminddb.Open(file)
	if err != nil {
		return nil, err
	}

	supportDefaultLang := false
	supportEnglish := false
	for _, lang := range db.Metadata.Languages {
		if lang == Lang {
			supportDefaultLang = true
			break
		}
		if lang == "en" {
			supportEnglish = true
		}
	}
	if !supportDefaultLang && supportEnglish {
		Lang = "en"
	}

	meta := model.Meta{
		Fields: FullFields,
	}

	switch db.Metadata.IPVersion {
	case 4:
		meta.IPVersion |= model.IPv4
	case 6:
		meta.IPVersion |= model.IPv4
		meta.IPVersion |= model.IPv6
	}

	return &Database{
		meta: meta,
		db:   db,
	}, nil
}

// Find 查询 IP 对应的网段和结果
func (d *Database) Find(ip net.IP) (*ipx.Range, map[string]string, error) {

	var data City
	ipNet, _, err := d.db.LookupNetwork(ip, &data)
	if err != nil {
		return nil, nil, err
	}

	return ipx.NewRange(ipNet), FieldsFormat(data.Format()), nil
}

// Meta 返回元数据
func (d *Database) Meta() model.Meta {
	return d.meta
}

// Close 关闭数据库实例
func (d *Database) Close() error {
	return d.db.Close()
}
