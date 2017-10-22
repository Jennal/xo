package internal

import "testing"
import "github.com/stretchr/testify/assert"
import "github.com/davecgh/go-spew/spew"

func TestParseExtraData(t *testing.T) {
	data, err := parseExtraData("`xo:name=ExName,type=string` 新增ExName字段")
	spew.Dump(data)
	assert.NotNil(t, data)
	assert.Nil(t, err)
}
