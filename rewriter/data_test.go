package rewriter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var DataLoad1 = `isp	Beihai	电信
province	上海市	上海
province	云南省	云南
[wrong line]
asnumber	4808	isp|联通
`

func TestDataRewriter(t *testing.T) {
	ast := assert.New(t)

	dataLoader := NewDataLoader()
	dataLoader.LoadString(DataLoad1)
	dataRewriter := NewDataRewriter(dataLoader, nil)
	data := map[string]string{}
	_, _, retData, err := dataRewriter.Rewrite(nil, nil, data)
	ast.Nil(err)
	ast.Equal(0, len(retData))

	data = map[string]string{
		"province": "上海市",
		"isp":      "Beihai",
	}
	_, _, retData, err = dataRewriter.Rewrite(nil, nil, data)
	ast.Nil(err)
	ast.Equal(2, len(retData))
	ast.Equal("上海", retData["province"])
	ast.Equal("电信", retData["isp"])

	data = map[string]string{
		"asnumber": "4808",
	}

	_, _, retData, err = dataRewriter.Rewrite(nil, nil, data)
	ast.Nil(err)
	ast.Equal(2, len(retData))
	ast.Equal("4808", retData["asnumber"])
	ast.Equal("联通", retData["isp"])
}
