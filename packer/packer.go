package packer

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"
)

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

const (
	FieldsSep        = "\t"
	IPv4      uint16 = 0x01
	IPv6      uint16 = 0x02
)

// Packer IPDB 打包工具
type Packer struct {
	Meta      Meta
	node      [][2]int
	dataHash  map[string]int
	dataChunk *bytes.Buffer
}

// Meta 元数据
type Meta struct {

	// Build 构建时间 10位时间戳
	Build int `json:"build"`

	// IPVersion IP库版本
	// IPv4: 0x01 IPv6: 0x02
	IPVersion uint16 `json:"ip_version"`

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

// NewPacker 初始化打包工具
func NewPacker(ipVersion uint16, languages map[string]int, fields []string) *Packer {
	if languages == nil || len(languages) == 0 {
		languages = map[string]int{"CN": 0}
	}
	return &Packer{
		Meta: Meta{
			Build:     int(time.Now().Unix()),
			IPVersion: ipVersion,
			Languages: languages,
			Fields:    fields,
		},
		node:      [][2]int{{}},
		dataChunk: &bytes.Buffer{},
		dataHash:  make(map[string]int),
	}
}

// Load 加载数据
func (p *Packer) Load(cidr string, fields []string) int {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return len(p.node)
	}
	if len(fields) != len(p.Meta.Fields) {
		return len(p.node)
	}
	mask, _ := ipNet.Mask.Size()

	node, index, ok := p.Nodes(ip, mask)
	if !ok {
		log.Printf("load cidr failed cidr(%s) data(%s) node(%d) index(%d) preview data(%s)\n", cidr, fields, node, index, p.resolve(-p.node[node][index]))
		return len(p.node)
	}
	if p.node[node][index] > 0 {
		log.Printf("cidr conflict %s %s\n", cidr, fields)
		return len(p.node)
	}
	offset := p.Fields(fields)
	p.node[node][index] = -offset
	return len(p.node)
}

func (p *Packer) resolve(offset int) string {
	offset -= 8
	data := p.dataChunk.Bytes()
	size := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	if (offset + 2 + size) > len(data) {
		return ""
	}
	return string(data[offset+2 : offset+2+size])
}

// Export 导出数据
func (p *Packer) Export() []byte {

	// Node Chunk
	nodeChunk := &bytes.Buffer{}
	p.Meta.NodeCount = len(p.node)
	for i := 0; i < p.Meta.NodeCount; i++ {
		for j := 0; j < 2; j++ {
			// 小于0: 数据记录，设置为NodeLength+Data偏移量
			// 等于0: 空值，设置为NodeLength
			// 大于0: Node跳转记录，不做调整
			if p.node[i][j] <= 0 {
				p.node[i][j] = p.Meta.NodeCount - p.node[i][j]
			}
			nodeChunk.Write(IntToBinaryBE(p.node[i][j], 32))
		}
	}
	// loopBack node
	nodeChunk.Write(IntToBinaryBE(p.Meta.NodeCount, 32))
	nodeChunk.Write(IntToBinaryBE(p.Meta.NodeCount, 32))

	// MetaData Chunk
	metaDataChunk := &bytes.Buffer{}
	p.Meta.TotalSize = nodeChunk.Len() + p.dataChunk.Len()
	metaData, err := json.Marshal(p.Meta)
	if err != nil {
		return nil
	}
	metaDataChunk.Write(IntToBinaryBE(len(metaData), 32))
	metaDataChunk.Write(metaData)

	// Result
	result := &bytes.Buffer{}
	if _, err := metaDataChunk.WriteTo(result); err != nil {
		return nil
	}
	if _, err := nodeChunk.WriteTo(result); err != nil {
		return nil
	}
	if _, err := p.dataChunk.WriteTo(result); err != nil {
		return nil
	}

	return result.Bytes()
}

// Nodes 获取CIDR地址所在节点和index
// 将补全Node中间链路，如果中间链路已经有数据，将无法写入新数据
func (p *Packer) Nodes(ip net.IP, mask int) (node, index int, ok bool) {
	// 如果传入的是IPv4，子网掩码增加96位( len(IPv6)-len(IPv4) )
	// 统一扩展为IPv6的子网掩码进行处理
	maxMask := mask - 1
	if ip.To4() != nil {
		maxMask += 96
	}
	for i := 0; i < maxMask; i++ {
		index = ((0xFF & int(ip[i>>3])) >> uint(7-(i%8))) & 1
		if p.node[node][index] == 0 {
			p.node = append(p.node, [2]int{})
			p.node[node][index] = len(p.node) - 1
		}
		if p.node[node][index] < 0 {
			return node, index, false
		}
		node = p.node[node][index]
	}
	return node, ((0xFF & int(ip[maxMask>>3])) >> uint(7-(maxMask%8))) & 1, true
}

// Fields 保存数据并返回数据的偏移量
// 相同的数据仅保存一份
// 数据格式 2 byte length + n byte data
func (p *Packer) Fields(fields []string) int {
	data := strings.Join(fields, FieldsSep)
	if _, ok := p.dataHash[data]; !ok {
		_data := []byte(data)
		// +8 是由于 loopBack node 占用了8byte
		p.dataHash[data] = p.dataChunk.Len() + 8
		p.dataChunk.Write(IntToBinaryBE(len(_data), 16))
		p.dataChunk.Write(_data)
	}
	return p.dataHash[data]
}

// IntToBinaryBE 将int转换为 binary big endian
func IntToBinaryBE(num, length int) []byte {
	switch length {
	case 16:
		_num := uint16(num)
		return []byte{byte(_num >> 8), byte(_num)}
	case 32:
		_num := uint32(num)
		return []byte{byte(_num >> 24), byte(_num >> 16), byte(_num >> 8), byte(_num)}
	default:
		return []byte{}
	}
}
