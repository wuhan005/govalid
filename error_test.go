package govalid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrContext_SetMessage(t *testing.T) {
	ruleCtx := ruleContext{
		field:  nil,
		rule:   "testrule",
		params: []string{"1", "a", "!"},
		value:  "testtest",
	}
	errCtx := NewErrorContext(ruleCtx)
	errCtx.SetMessage("TestMessage")

	assert.Equal(t, errCtx.Message, "TestMessage")
}

func TestGetErrorTemplate(t *testing.T) {
	assert.Equal(t, GetErrorTemplate("not_found_name"), "未知错误")
	assert.Equal(t, GetErrorTemplate("required"), "不能为空")
}
