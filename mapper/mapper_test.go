package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestProvinceMapFile = "../data/province.map"
)

func TestDataMapper(t *testing.T) {
	ast := assert.New(t)

	mapper := NewDataMapper(TestProvinceMapFile)
	ast.NotNil(mapper)

	replace, matched := mapper.Mapping("province", "上海市")
	ast.True(matched)
	ast.Equal("上海", replace)
}
