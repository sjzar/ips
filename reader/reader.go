package reader

import (
	"net"
	"strings"

	"github.com/sjzar/ips/db"
	"github.com/sjzar/ips/iprange"
	"github.com/sjzar/ips/mapper"
	"github.com/sjzar/ips/model"
)

const FieldsSep = ":"

// FilterReader 支持字段过滤的IP库读取器
type FilterReader struct {
	// database IP库
	database db.Database

	// filter 字段过滤
	filter Filter

	// mapper 数据映射
	mapper mapper.Mapper

	// Calibrate 数据校验
	// 用于做数据对比与校准
	Calibrate func(ip net.IP, data map[string]string) (ipNet *net.IPNet, ret map[string]string, updated bool)
}

// NewFilterReader 初始化IP库读取器
func NewFilterReader(database db.Database, filter Filter, mapper mapper.Mapper) *FilterReader {
	return &FilterReader{
		database: database,
		filter:   filter,
		mapper:   mapper,
	}
}

// FieldsCollection 字段集合
func (r *FilterReader) FieldsCollection() string {
	return r.filter.FieldsCollection()
}

// LookupNetwork 查询IP所在网段和对应数据
func (r *FilterReader) LookupNetwork(ip net.IP) (*net.IPNet, string, error) {
	ipNet, data, err := r.database.LookupNetwork(ip)
	if err != nil {
		return nil, "", err
	}

	// Data Mapping
	for k := range data {
		if _val, matched := r.mapper.Mapping(k, data[k]); matched {
			data[k] = _val
		}
	}

	if r.Calibrate != nil {
		if _ipNet, _data, updated := r.Calibrate(ip, data); updated {
			if !iprange.IPNetMaskLess(ipNet, _ipNet) {
				ipNet = _ipNet
			}
			data = _data
		}
	}

	fields := r.filter.Fields(data[r.filter.Key()])
	ret := make([]string, len(fields))
	for i, field := range fields {
		if field == model.Placeholder {
			continue
		}
		ret[i] = data[field]
	}
	return ipNet, strings.Join(ret, FieldsSep), nil
}
