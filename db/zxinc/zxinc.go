package zxinc

import (
	"net"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

type Database struct {
	meta model.Meta
	db   *ZXInc
}

// New 初始化 QQWRY 数据库实例
func New(file string) (*Database, error) {

	db, err := NewZXInc(file)
	if err != nil {
		return nil, err
	}

	meta := model.Meta{
		Fields:    []string{FieldCountry, FieldArea},
		IPVersion: model.IPv6,
	}

	return &Database{
		meta: meta,
		db:   db,
	}, nil
}

// Find 查询 IP 对应的网段和结果
func (d *Database) Find(ip net.IP) (*ipx.Range, map[string]string, error) {
	ipr, data, err := d.db.Find(ip)
	if err != nil {
		return nil, nil, err
	}
	return ipr, FieldsFormat(data), err
}

// Meta 返回元数据
func (d *Database) Meta() model.Meta {
	return d.meta
}

// Close 关闭数据库实例
func (d *Database) Close() error {
	return nil
}
