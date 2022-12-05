package ipio

import (
	"io"
	"net"

	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

// Reader IP库读取工具
// 从数据库中获取IP的地理位置信息，并格式化数据
type Reader interface {

	// Meta 返回元数据
	Meta() model.Meta

	// Find 查询IP所在网段和对应数据
	Find(ip net.IP) (*ipx.Range, []string, error)

	// Close 关闭
	Close() error
}

// Scanner IP库扫描工具
type Scanner interface {

	// Meta 返回元数据
	Meta() model.Meta

	// Scan 扫描下一条数据，返回 true 表示有数据，false 表示没有数据
	Scan() bool

	// Result 返回结果
	Result() (*ipx.Range, []string)

	// Err 返回错误
	Err() error

	// Close 关闭
	Close() error
}

// Writer IP库写入工具
type Writer interface {

	// Insert 插入数据
	Insert(ipr *ipx.Range, values []string) error

	// Save 保存数据
	Save(w io.Writer) error
}
