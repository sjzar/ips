package zxinc

import (
	"github.com/sjzar/ips/model"
)

// cidr,country,area
// 2001:200:120::/48,日本    东京都  品川区,Sony Computer Science 研究所
// 2001:200:178::/48,美国    California州    San José,Sony NCSA实验室
// 2001:250:1::/48,中国    北京市,教育网(CERNET)网络运行部
// ffff:ffff:ffff:fff0::/60,,ZX公网IPv6库 20210511版
// * Country中有多级数据，需要处理

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
