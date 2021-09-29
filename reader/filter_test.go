package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjzar/ips/model"
)

func TestFilter(t *testing.T) {
	ast := assert.New(t)

	// new filter with magicArg
	filter := NewKeyFieldsFilter("")
	ast.NotNil(filter)
	ast.Equal(model.Country, filter.Key())
	ast.Equal(4, len(filter.Fields("")))
	ast.Equal("country,province,city,isp", filter.FieldsCollection())

	arg := "city#温州:province,city,isp|province,-,isp"
	filter = NewKeyFieldsFilter(arg)
	ast.Equal(model.City, filter.Key())
	ast.Equal(model.City, filter.Fields("温州")[1])
	ast.Equal(model.Placeholder, filter.Fields("上海")[1])
	ast.Equal("province,city,isp", filter.FieldsCollection())
}
