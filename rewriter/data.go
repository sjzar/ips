package rewriter

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
)

const (
	DataSep    = "\t"
	ReplaceSep = "|"
)

// DataRewriter 数据改写
type DataRewriter struct {
	DataLoader *DataLoader
	rw         Rewriter
}

// NewDataRewriter 初始化数据改写
// 数据以'\t'作为分段，每一行的格式为: <field>\t<match>\t<replace>\n
// @ <field> - 字段ID
// @ <match> - 匹配内容
// @ <replace> - 改写内容，支持以'|'作为分段，表示匹配后改写其他字段
// 举例:
// country\t\t保留地址 - "国家"字段中，如果数据为空，改写为"保留地址"
// province\t内蒙古\t内蒙 - "省份"字段中，如果数据为"内蒙古"，改写为"内蒙"
// asnumber\t4134\tisp|电信 - "AS号码"字段中，如果数据为"4134"，改写"运营商"字段为"电信"
func NewDataRewriter(dl *DataLoader, rw Rewriter) *DataRewriter {
	if rw == nil {
		rw = DefaultRewriter
	}

	if dl == nil {
		dl = NewDataLoader()
	}

	return &DataRewriter{
		DataLoader: dl,
		rw:         rw,
	}
}

// Rewrite 改写
func (r *DataRewriter) Rewrite(ip net.IP, ipRange *ipx.Range, data map[string]string) (net.IP, *ipx.Range, map[string]string, error) {
	r.DataLoader.Rewrite(data)
	return r.rw.Rewrite(ip, ipRange, data)
}

// DataLoader 数据加载
type DataLoader struct {
	// field, match, replace
	Data []map[string]map[string]string
}

// NewDataLoader 初始化数据加载
func NewDataLoader() *DataLoader {
	return &DataLoader{
		Data: make([]map[string]map[string]string, 0),
	}
}

// LoadFile 加载文件
func (l *DataLoader) LoadFile(file string) error {
	if len(file) == 0 {
		return errors.ErrEmptyFile
	}
	f, err := os.Open(file)
	if err != nil {
		log.Println("open file failed", err)
		return err
	}
	defer f.Close()

	return l.load(f)
}

// LoadString 加载数据
func (l *DataLoader) LoadString(data ...string) {
	l.load(strings.NewReader(strings.Join(data, "\n")))
}

// load 加载数据
func (l *DataLoader) load(r io.Reader) error {
	if l.Data == nil {
		l.Data = make([]map[string]map[string]string, 0)
	}
	data := make(map[string]map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		split := strings.SplitN(string(line), DataSep, 3)
		if len(split) < 3 {
			continue
		}
		if _, ok := data[split[0]]; !ok {
			data[split[0]] = make(map[string]string)
		}
		data[split[0]][split[1]] = split[2]
	}
	if err := scanner.Err(); err != nil {
		log.Println("load data failed", err)
		return err
	}
	l.Data = append(l.Data, data)
	return nil
}

// Rewrite 改写
func (l *DataLoader) Rewrite(data map[string]string) {
	if l.Data == nil || len(l.Data) == 0 {
		return
	}
	for _, _data := range l.Data {
		for field, match := range _data {
			value, ok := data[field]
			if !ok {
				continue
			}
			replace, ok := match[value]
			if !ok {
				continue
			}
			split := strings.Split(replace, ReplaceSep)
			if len(split) > 1 {
				data[split[0]] = split[1]
			} else {
				data[field] = split[0]
			}
		}
	}
}
