package db

import (
	"net"
	"path/filepath"

	"github.com/sjzar/ips/db/awdb"
	"github.com/sjzar/ips/db/ipdb"
	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

const (
	// FormatIPDB IPDB格式
	// Official: https://www.ipip.net/
	FormatIPDB = "ipdb"

	// FormatAWDB AWDB格式
	// Official: https://www.ipplus360.com/
	FormatAWDB = "awdb"
)

// Database IP 数据库
type Database interface {

	// Meta 返回元数据
	Meta() model.Meta

	// Find 查询 IP 对应的网段和结果
	Find(ip net.IP) (*ipx.Range, map[string]string, error)

	// Close 关闭数据库实例
	Close() error
}

// NewDatabase 初始化 IP 数据库
func NewDatabase(format, file string) (Database, error) {
	if len(format) == 0 {
		switch {
		case filepath.Ext(file) == ".ipdb":
			format = FormatIPDB
		case filepath.Ext(file) == ".awdb":
			format = FormatAWDB
		}
	}

	switch format {
	case FormatIPDB:
		return ipdb.New(file)
	case FormatAWDB:
		return awdb.New(file)
	}

	return nil, errors.ErrDBFormatNotSupported
}
