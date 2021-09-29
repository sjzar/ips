package reader

import "net"

// Reader IP库读取工具
// 从数据库中获取IP的地理位置信息，并格式化数据
type Reader interface {

	// LookupNetwork 查询IP所在网段和对应数据
	LookupNetwork(ip net.IP) (*net.IPNet, string, error)

	// FieldsCollection 字段集合
	FieldsCollection() string
}

// Filter 字段过滤工具
type Filter interface {

	// Key 条件字段
	Key() string

	// Fields 查询字段列表
	// 根据条件字段对应的值不同，输出不同的字段
	Fields(key string) []string

	// FieldsCollection 字段集合
	// 返回过滤器完整的字段列表
	FieldsCollection() string
}
