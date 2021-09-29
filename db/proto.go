package db

import "net"

// Database IP数据库的查询封装
type Database interface {
	LookupNetwork(ip net.IP) (*net.IPNet, map[string]string, error)
}
