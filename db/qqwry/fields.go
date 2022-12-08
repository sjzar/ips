package qqwry

import (
	"github.com/sjzar/ips/model"
)

// cidr,country,area
// 1.2.2.0/24,北京市海淀区,北龙中网(北京)科技有限公司
// 1.2.8.0/24,中国,网络信息中心
// 1.2.9.0/24,广东省,电信
// 1.8.151.0/24,香港,特别行政区
// 1.10.140.0/25,泰国, CZ88.NET
// 1.12.36.0/22,广东省广州市,腾讯云
// * Country 和 Area 中的数据比较混乱，需要后续改写处理

const (
	FieldCountry = "country"
	FieldArea    = "area"
)

// FullFields 全字段列表
var FullFields = []string{
	FieldCountry,
	FieldArea,
}

// CommonFieldsMap 公共字段映射
var CommonFieldsMap = map[string]string{
	model.Country: FieldCountry,
	model.ISP:     FieldArea,
}

// FieldsFormat 字段格式化，并补充公共字段
func FieldsFormat(data map[string]string) map[string]string {

	// Fill Common Fields
	for k, v := range CommonFieldsMap {
		data[k] = data[v]
	}

	return data
}
