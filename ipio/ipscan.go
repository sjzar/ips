package ipio

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"strings"

	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

/*
定义 IPScan 文件，作为 IP 数据库扫描后平铺的文件格式，用于检查和修改 IP 数据库
1. 实现 Writer 接口，用于写入 IPScan 文件，可与 DBScanner 配合使用，扫描 IP 数据库并写入 IPScan 文件
2. 实现 Scanner 接口，用于读取 IPScan 文件，可与 IPDB Writer 一起使用，实现 IPDB 的打包

IPScan 文件格式：
# ScanTime: 2022-12-04 16:00:00
# IPVersion: 4
# Fields: country,province,city,isp
# Meta: {}
0.0.0.0/0	国家,省份,城市,运营商
*/

const (
	CommentPrefix = "#"
	MetaPrefix    = "# Meta: "
	FieldSep      = ","
)

// IPScanScanner IPScan 扫描工具
type IPScanScanner struct {
	FilePath string
	meta     model.Meta
	file     *os.File
	scanner  *bufio.Scanner

	ipr    *ipx.Range
	values []string
	done   bool
	err    error
}

// NewIPScanScanner 初始化 IPScan 扫描实例
func NewIPScanScanner(filePath string) (*IPScanScanner, error) {

	s := &IPScanScanner{
		FilePath: filePath,
	}

	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

// init 初始化 IPScan 扫描实例
func (s *IPScanScanner) init() error {

	var err error
	s.file, err = os.Open(s.FilePath)
	if err != nil {
		return err
	}

	metaflag := false
	s.scanner = bufio.NewScanner(s.file)
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, MetaPrefix) {
			metaflag = true
			metaStr := line[len(MetaPrefix):]
			if err := json.Unmarshal([]byte(metaStr), &s.meta); err != nil {
				return err
			}
			break
		}

		if strings.HasPrefix(line, CommentPrefix) {
			continue
		}
	}

	if !metaflag {
		return errors.ErrMetaNotFound
	}

	return nil
}

// Meta 返回元数据
func (s *IPScanScanner) Meta() model.Meta {
	return s.meta
}

// Scan 扫描下一条数据，返回 true 表示有数据，false 表示没有数据
func (s *IPScanScanner) Scan() bool {
	if s.done {
		return false
	}
	s.ipr, s.values = nil, nil

	for s.scanner.Scan() {
		line := s.scanner.Text()
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		split := strings.SplitN(line, "\t", 2)
		if len(split) != 2 {
			continue
		}

		_, ipNet, err := net.ParseCIDR(split[0])
		if err != nil {
			continue
		}

		values := strings.Split(split[1], FieldSep)
		if len(values) != len(s.meta.Fields) {
			continue
		}

		s.ipr, s.values = ipx.NewRange(ipNet), values
		if s.ipr.IsEnd() {
			s.done = true
		}
		break
	}

	if err := s.scanner.Err(); err != nil {
		s.err = err
		return false
	}

	return true
}

// Result 返回结果
func (s *IPScanScanner) Result() (*ipx.Range, []string) {
	return s.ipr, s.values
}

// Err 返回错误
func (s *IPScanScanner) Err() error {
	return s.err
}

// Close 关闭
func (s *IPScanScanner) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
