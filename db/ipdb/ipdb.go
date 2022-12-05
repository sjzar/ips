package ipdb

import (
	"net"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

// Database IPDB 数据库
type Database struct {
	meta model.Meta
	db   *City
}

// New 初始化 IPDB 数据库实例
func New(file string) (*Database, error) {
	city, err := NewCity(file)
	if err != nil {
		return nil, err
	}

	meta := model.Meta{
		Fields: city.Fields(),
	}
	if city.IsIPv4() {
		meta.IPVersion |= model.IPv4
	}
	if city.IsIPv6() {
		meta.IPVersion |= model.IPv6
	}

	return &Database{
		meta: meta,
		db:   city,
	}, nil
}

// Find 查询 IP 对应的网段和结果
func (d *Database) Find(ip net.IP) (*ipx.Range, map[string]string, error) {
	data, ipNet, err := d.db.FindMap(ip.String(), "CN")
	if err != nil {
		return nil, nil, err
	}

	return ipx.NewRange(ipNet), FieldsFormat(data), nil
}

// Meta 返回元数据
func (d *Database) Meta() model.Meta {
	return d.meta
}

// Close 关闭数据库实例
func (d *Database) Close() error {
	return nil
}
