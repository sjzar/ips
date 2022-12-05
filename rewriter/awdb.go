package rewriter

import (
	"log"
	"net"

	"github.com/sjzar/ips/db/awdb"
	"github.com/sjzar/ips/ipx"
	"github.com/sjzar/ips/model"
)

const (
	AWDBFieldLine = "line"
)

// AWDBIPLineRewriter AWDB数据库的IP线路改写
// 根据线路数据库匹配 ISP
type AWDBIPLineRewriter struct {
	awdb *awdb.Database
	rw   Rewriter
}

// NewIPLineRewriter 初始化实例
func NewIPLineRewriter(ipLineFile string, rw Rewriter) (Rewriter, error) {
	if rw == nil {
		rw = DefaultRewriter
	}

	db, err := awdb.New(ipLineFile)
	if err != nil {
		log.Println("new awdb database failed", err)
		return nil, err
	}
	return &AWDBIPLineRewriter{
		awdb: db,
		rw:   rw,
	}, nil
}

// Rewrite 改写
func (r *AWDBIPLineRewriter) Rewrite(ip net.IP, ipRange *ipx.Range, data map[string]string) (net.IP, *ipx.Range, map[string]string, error) {
	ipRangeLine, dataLine, err := r.awdb.Find(ip)
	if err != nil {
		return nil, nil, nil, err
	}

	if line, ok := dataLine[AWDBFieldLine]; ok {
		if ok := ipRange.CommonRange(ip, ipRangeLine); ok {
			data[model.ISP] = line
		}
	}

	return r.rw.Rewrite(ip, ipRange, data)
}
