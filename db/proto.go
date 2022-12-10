package db

import (
	"net"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sjzar/ips/db/awdb"
	"github.com/sjzar/ips/db/ipdb"
	"github.com/sjzar/ips/db/mmdb"
	"github.com/sjzar/ips/db/qqwry"
	"github.com/sjzar/ips/db/zxinc"
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

	// FormatMMDB MMDB格式
	// Official: https://www.maxmind.com/
	FormatMMDB = "mmdb"

	// FormatQQWry QQWry格式
	// Official: https://www.cz88.net/
	FormatQQWry = "qqwry"

	// FormatZXInc ZXInc格式
	// Official: https://ip.zxinc.org/
	FormatZXInc = "zxinc"
)

const (
	// FormatSep 格式分隔符
	FormatSep = ":"
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

	// support format:file
	// ipdb:./data/ipip.ipdb
	// awdb:./data/ipplus.awdb
	if len(format) == 0 && strings.Contains(file, FormatSep) {
		split := strings.SplitN(file, FormatSep, 2)
		format = split[0]
		file = split[1]
		if strings.Contains(file, "~") {
			if u, err := user.Current(); err == nil {
				file = strings.Replace(file, "~", u.HomeDir, -1)
			}
		}
		if _file, err := filepath.Abs(file); err != nil {
			file = _file
		}
	}

	if len(format) == 0 {
		switch {
		case filepath.Ext(file) == ".ipdb":
			format = FormatIPDB
		case filepath.Ext(file) == ".awdb":
			format = FormatAWDB
		case filepath.Ext(file) == ".mmdb":
			format = FormatMMDB
		case strings.HasSuffix(file, "qqwry.dat"):
			format = FormatQQWry
		case strings.HasSuffix(file, "zxipv6wry.db"):
			format = FormatZXInc
		}
	}

	switch format {
	case FormatIPDB:
		return ipdb.New(file)
	case FormatAWDB:
		return awdb.New(file)
	case FormatQQWry:
		return qqwry.New(file)
	case FormatMMDB:
		return mmdb.New(file)
	case FormatZXInc:
		return zxinc.New(file)
	}

	return nil, errors.ErrDBFormatNotSupported
}
