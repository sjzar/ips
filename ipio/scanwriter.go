package ipio

// ScanWriter IP库扫描写入工具
type ScanWriter struct {
	Scanner
	Writer
}

// NewScanWriter 初始化 IP 库扫描写入实例
func NewScanWriter(s Scanner, w Writer) *ScanWriter {
	return &ScanWriter{
		Scanner: s,
		Writer:  w,
	}
}

// ScanWrite 扫描并写入
func (sw *ScanWriter) ScanWrite() error {
	for sw.Scan() {
		if err := sw.Writer.Insert(sw.Result()); err != nil {
			return err
		}
	}
	if err := sw.Err(); err != nil {
		return err
	}
	return nil
}

// ScanWrite 扫描并写入
func ScanWrite(s Scanner, w Writer) error {
	return NewScanWriter(s, w).ScanWrite()
}
