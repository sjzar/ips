package mapper

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

// DataMapper 数据映射工具
type DataMapper struct {
	// field, match, replace
	data map[string]map[string]string
}

// Mapping 字段映射
func (m *DataMapper) Mapping(field, match string) (string, bool) {
	if m.data == nil {
		return "", false
	}
	if _map, ok := m.data[field]; ok {
		if data, ok := _map[match]; ok {
			return data, true
		}
	}
	return "", false
}

// NewDataMapper 初始化数据映射工具
func NewDataMapper(file string) *DataMapper {
	mapper := &DataMapper{
		data: make(map[string]map[string]string),
	}

	if len(file) != 0 {
		mapper.LoadFile(file)
	}

	return mapper
}

// LoadFile 加载映射文件
func (m *DataMapper) LoadFile(file string) {
	if len(file) == 0 {
		return
	}
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if m.data == nil {
		m.data = make(map[string]map[string]string)
	}

	m.load(f)
}

// LoadString 加载映射数据
func (m *DataMapper) LoadString(data string) {
	m.load(strings.NewReader(data))
}

// load 加载映射数据
// 映射数据以\t作为分段
// <field>\t<原始数据>\t<修正数据>
// 例如:
// country\t\t保留地址
// province\t内蒙古\t内蒙
// * 兼容忽略field的情况
func (m *DataMapper) load(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		split := strings.SplitN(string(line), "\t", 3)
		if len(split) < 2 {
			continue
		}
		key := ""
		if len(split) == 3 {
			key = split[0]
			split[0], split[1] = split[1], split[2]
		}

		if _, ok := m.data[key]; !ok {
			m.data[key] = make(map[string]string)
		}

		m.data[key][split[0]] = split[1]
	}
	if err := scanner.Err(); err != nil {
		log.Print(err)
	}
}
