/*
 * Copyright (c) 2023 shenjunzheng@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ipdb

// Meta ipdb 元数据
type Meta struct {

	// Build 构建时间 10位时间戳
	Build int `json:"build"`

	// IPVersion IP库版本
	IPVersion int `json:"ip_version"`

	// Languages 支持语言
	// value为语言对应的fields偏移量
	Languages map[string]int `json:"languages"`

	// NodeCount 节点数量
	NodeCount int `json:"node_count"`

	// TotalSize 节点区域和数据区域大小统计
	TotalSize int `json:"total_size"`

	// Fields 数据字段列表
	// 城市级别数据库包含13个字段
	// "country_name": "国家名称"
	// "region_name": "省份名称"
	// "city_name": "城市名称"
	// "owner_domain": "所有者"
	// "isp_domain": "运营商"
	// "latitude": "纬度"
	// "longitude": "经度"
	// "timezone": "时区"
	// "utc_offset": "UTC偏移量"
	// "china_admin_code": "中国邮编"
	// "idd_code": "电话区号"
	// "country_code": "国家代码"
	// "continent_code": "大陆代码"
	Fields []string `json:"fields"`
}
