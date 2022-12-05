package awdb

import (
	"net"

	"github.com/dilfish/awdb-golang/awdb-golang"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

// Database AWDB 数据库
type Database struct {
	meta model.Meta
	db   *awdb.Reader
}

// New 初始化 AWDB 数据库实例
func New(file string) (*Database, error) {
	db, err := awdb.Open(file)
	if err != nil {
		return nil, err
	}

	meta := model.Meta{
		Fields: FullFields,
	}
	if db.Metadata.IPVersion == 4 {
		meta.IPVersion |= model.IPv4
	} else if db.Metadata.IPVersion == 6 {
		meta.IPVersion |= model.IPv6
	}

	return &Database{
		meta: meta,
		db:   db,
	}, nil
}

// Meta 返回元数据
func (d *Database) Meta() model.Meta {
	return d.meta
}

// Find 查询 IP 对应的网段和结果
func (d *Database) Find(ip net.IP) (*ipx.Range, map[string]string, error) {
	var record interface{}
	ipNet, ok, err := d.db.LookupNetwork(ip, &record)
	if !ok || err != nil {
		return nil, nil, err
	}

	return ipx.NewRange(ipNet), FieldsFormat(record.(map[string]interface{})), nil
}

// Close 关闭数据库实例
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}
