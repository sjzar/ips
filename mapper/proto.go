package mapper

// Mapper 数据映射工具
// 替换原始数据，用于统一输出数据格式
type Mapper interface {

	// Mapping 字段映射
	Mapping(field, match string) (string, bool)
}
