package zxinc

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"os"

	"github.com/sjzar/ips/errors"
	"github.com/sjzar/ips/ipx"
)

const (
	// RedirectMode1 重定向模式1
	// 表示国家记录和地区记录都被重定向
	RedirectMode1 = 0x01

	// RedirectMode2 重定向模式2
	// 表示国家记录或地区记录被重定向
	RedirectMode2 = 0x02
)

/* zxipv6wry.dat Format (Little Endian)
+--------------------------------+--------------------------------+
|         "IPDB" (4byte)         |         Version (2byte)        |
+--------------------------------+--------------------------------+
|      Offset Length (1byte)     |        IP Length (1byte)       |
+--------------------------------+--------------------------------+
|         Count (8byte)          |        Start Index (8byte)     |
+--------------------------------+--------------------------------+
|                            Data Chunk                           |
+--------------------------------+--------------------------------+
|                            Index Chunk                          |
+--------------------------------+--------------------------------+

Data Chunk like qqwry.dat, but the End IP not included and use UTF-8 encoding.

*/

// ZXInc ZXInc 数据库
type ZXInc struct {

	// data IP库数据
	data []byte

	// start IP库数据开始位置
	start uint64

	// end IP库数据结束位置
	end uint64

	// version IP库版本, 一般是 0x1
	version []byte

	// offsetLen 偏移地址长度，一般是 3
	offsetLen uint64

	// ipLen IP地址长度，(4,8,12,16)，目前是 8
	ipLen uint64

	// indexLen 索引长度, ipLen + offsetLen 为一条索引
	indexLen uint64
}

func NewZXInc(filePath string) (*ZXInc, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) < 24 {
		return nil, errors.ErrDatabaseIsInvalid
	}

	version := data[4:5]
	offsetLen := uint64(data[6])
	ipLen := uint64(data[7])
	count := binary.LittleEndian.Uint64(data[8:16])
	start := binary.LittleEndian.Uint64(data[16:24])
	indexLen := offsetLen + ipLen
	end := start + count*indexLen

	if uint64(len(data)) < end || start >= end {
		return nil, errors.ErrDatabaseIsInvalid
	}

	return &ZXInc{
		data:      data,
		start:     start,
		end:       end,
		version:   version,
		offsetLen: offsetLen,
		ipLen:     ipLen,
		indexLen:  indexLen,
	}, nil
}

// Find 查找IP
func (q *ZXInc) Find(ip net.IP) (*ipx.Range, map[string]string, error) {

	ip = ip.To16()
	if ip == nil {
		return nil, nil, errors.ErrIPVersionNotSupported
	}

	startIP, nextIP, offset := q.findOffset(binary.BigEndian.Uint64(ip[:q.ipLen]))
	if offset == 0 {
		return nil, nil, errors.ErrDatabaseIsInvalid
	}
	country, area, err := q.parse(offset, 0)
	if err != nil {
		return nil, nil, err
	}

	return &ipx.Range{
			Start: ipx.Uint64ToIP(startIP),
			End:   ipx.PrevIP(ipx.Uint64ToIP(nextIP)),
		}, map[string]string{
			FieldCountry: country,
			FieldArea:    area,
		}, nil
}

// findOffset 查找IP对应的偏移量
func (q *ZXInc) findOffset(ip uint64) (startIP, nextIP uint64, offset uint32) {
	low := q.start
	high := q.end

	var mid, currentIP uint64

	for {
		mid = (high-low)/q.indexLen/2*q.indexLen + low

		if high-low == q.indexLen {
			if binary.LittleEndian.Uint64(q.data[high:high+q.ipLen]) <= ip {
				mid = high
			}
			next := mid + q.indexLen
			if next > q.end {
				next = q.end
			}

			return binary.LittleEndian.Uint64(q.data[mid : mid+q.ipLen]),
				binary.LittleEndian.Uint64(q.data[next : next+q.ipLen]),
				Bytes3Uint32(q.data[mid+q.ipLen : mid+q.indexLen])
		}

		currentIP = binary.LittleEndian.Uint64(q.data[mid : mid+q.ipLen])
		if currentIP > ip {
			high = mid
		} else if currentIP < ip {
			low = mid
		} else {
			next := mid + q.indexLen
			if next > q.end {
				next = q.end
			}
			return ip, binary.LittleEndian.Uint64(q.data[next : next+q.ipLen]), Bytes3Uint32(q.data[mid+q.ipLen : mid+q.indexLen])
		}
	}
}

// parse 解析数据
func (q *ZXInc) parse(offset uint32, depth int) (country, area string, err error) {
	if depth > 1 {
		return "", "", errors.ErrDatabaseIsInvalid
	}

	switch q.data[offset] {
	case RedirectMode1:
		// Redirect Mode1: redirect country AND area
		return q.parse(Bytes3Uint32(q.data[offset+1:offset+4]), depth+1)
	case RedirectMode2:
		// Redirect Mode2: redirect country OR area
		country, _, err = q.parseString(Bytes3Uint32(q.data[offset+1 : offset+4]))
		if err != nil {
			return "", "", err
		}
		offset += 4
	default:
		var length int
		country, length, err = q.parseString(offset)
		if err != nil {
			return "", "", err
		}
		// +1 跳过结束标志(0x00)
		offset += uint32(length) + 1
	}
	area, err = q.parseArea(offset, depth)
	if err != nil {
		return "", "", err
	}

	return country, area, nil
}

// parseArea 解析地区
func (q *ZXInc) parseArea(offset uint32, depth int) (area string, err error) {
	if depth > 2 {
		return "", errors.ErrDatabaseIsInvalid
	}

	switch q.data[offset] {
	case RedirectMode1, RedirectMode2:
		return q.parseArea(Bytes3Uint32(q.data[offset+1:offset+4]), depth+1)
	}
	area, _, err = q.parseString(offset)
	if err != nil {
		return "", err
	}

	return area, nil
}

// parseString 解析字符串
func (q *ZXInc) parseString(offset uint32) (string, int, error) {
	length := bytes.IndexByte(q.data[offset:], 0x00)
	if length == -1 {
		return "", 0, errors.ErrDatabaseIsInvalid
	}
	str := string(q.data[offset : offset+uint32(length)])
	return str, length, nil
}

// Bytes3Uint32 3字节转换为uint32
func Bytes3Uint32(b []byte) uint32 {
	_ = b[2]
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}
