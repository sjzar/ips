package ipio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldSelector(t *testing.T) {
	ast := assert.New(t)

	selector := NewFieldSelector("chinaCity")
	ast.Equal([]string{"country", "province", "city", "isp"}, selector.Fields())

	data := map[string]string{
		"country":  "中国",
		"province": "浙江",
		"city":     "杭州",
		"isp":      "电信",
	}
	ast.Equal([]string{"中国", "浙江", "杭州", "电信"}, selector.Select(data))

	data = map[string]string{
		"country":  "日本",
		"province": "东京都",
		"city":     "品川区",
		"isp":      "WIDE Project",
	}
	ast.Equal([]string{"日本", "", "", ""}, selector.Select(data))
}
