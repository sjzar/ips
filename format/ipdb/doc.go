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

/* IPDB Format
	+--------------------------------+--------------------------------+
	|    MetaData Length (4byte)     |     MetaData (Json Format)     |
	+--------------------------------+--------------------------------+
	|                 Node Chunk (Prefix Tree / Trie)                 |
	+--------------------------------+--------------------------------+
	|                            Data Chunk                           |
	+--------------------------------+--------------------------------+

* 文件被划分为3块，分别是MetaData、NodeChunk、DataChunk
* MetaData是一个json，保存了数据库的构建信息和查询所需的一些元数据
* NodeChunk是前缀树(Trie)，每个Node为8byte，用于保存下一跳offset位置信息，若offset大于Node数量，将跳转到DataChunk，表示查询到了结果
* DataChunk用于保存IP库数据，相同的数据仅保存一份，减少数据冗余

# metadata 示例
{
	"build": 1632971142,    // 构建时间
	"ip_version": 1,        // IP库版本 IPv4:0x1 IPv6:0x2
	"languages": {
		"CN": 0             // 语言 & 字段偏移量
	},
	"node_count": 8705098,  // Node数量
	"total_size": 90028407, // NodeChunk+DataChunk 数据大小
	"fields": [             // DataChunk中每组数据的字段
		"country_name",
		"region_name",
		"city_name",
		"owner_domain",
		"isp_domain",
		"latitude",
		"longitude",
		"timezone",
		"utc_offset",
		"china_admin_code",
		"idd_code",
		"country_code",
		"continent_code"
	]
}

	+--------------------------------+--------------------------------+--------------------------------+
	| Data Length (2byte) | Data Fields(<country>\t<province>\t<city>\t<isp>\t<country>\t<province>)   |
	+--------------------------------+--------------------------------+--------------------------------+
	| Data Length (2byte) | Data Fields(<country>\t<province>\t<city>\t<isp>\t<country>\t<province>)   |
	+--------------------------------+--------------------------------+--------------------------------+
* DataChunk中，数据被分为两个部分，第一部分是长度，第二部分是数据
* 数据部分使用`\t`分隔字段，多语言版本的IP库采用字段offset返回不同语言的数据

* 查询
* CIDR地址是IP+子网掩码的一种网段的描述方式，例如10.0.0.1/8，表示子网掩码为8位（255.0.0.0），CIDR中所有的IP，子网掩码部分完全相同
* 将IP地址看作32bit/128bit的二进制字符串，在NodeChunk中使用前缀树(Trie)从前往后查询，匹配到CIDR后，跳转到DataChunk中返回对应数据
* CIDR分组越多，Node数量越大
* CIDR禁止互相包含，例如将 10.0.0.1/8 和 10.0.0.1/16 设置为不同的结果，在查询过程中，由于先匹配到了10.0.0.1/8，所以10.0.0.1/16 不生效，需要做拆分
* 详细的查询过程，可以参考这篇论文 IPv4 route lookup on Linux: https://vincent.bernat.ch/en/blog/2017-ipv4-route-lookup-linux

* 打包
* 构建前缀树，将不同的数据按照Load顺序放入DataChunk中
* 前缀树按照IPv6的规格进行构建，IPv4需要补全前面96bit的子网数据，80bit的0和16bit的1，对应IPv6的映射地址（::FFFF:<IPv4>）；查询IPv4时，支持快速offset到96bit mask位置起进行查询
* ipdb格式数据库理论上支持IPv4/IPv6数据处于同一个数据库中，但是需要注意::FFFF:路径上禁止存在其他CIDR记录，保证到IPv4查询的路径畅通（修改查询SDK可以更好的支持IPv4/IPv6同时查询）

*/
