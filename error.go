package govalid

import (
	"fmt"
	"strings"
)

var _ error = (*ErrContext)(nil)

// ErrContext contains the error context.
type ErrContext struct {
	FieldName  string
	FieldLabel string
	FieldValue interface{}

	fieldLimitValue interface{}
	errorTemplate   string
	errorMessage    string
}

func (e *ErrContext) Error() string {
	return e.errorMessage
}

var (
	FieldNamePlaceholder  = "{field}"
	FieldLimitPlaceholder = "{limit}"
)

var errorTemplate = map[string]string{
	"required":       "不能为空",
	"min":            "应大于",
	"max":            "应小于",
	"minlen":         "长度应大于",
	"maxlen":         "长度应小于",
	"alpha":          "必须只包含字母",
	"alphanumeric":   "只能含有字母或数字",
	"alphadash":      "只含有数字或字母以及下划线",
	"firstCharAlpha": "的第一个字符必须为字母",
	"lastUnderline":  "的最后一个字符不能为下划线",
	"email":          "不是合法的电子邮箱格式",
	"ipv4":           "不是合法的 IPv4 地址格式",
	"mobile":         "不是合法的手机号",
	"tel":            "不是合法的座机号码",
	"phone":          "不是合法的号码",
	"idcard":         "不是合法的身份证号",
	"equal":          "的值前后不相同",

	"_checkerNotFound":      "检查规则未找到}}",
	"_unknownErrorTemplate": "{{未知错误}}",
	"_paramError":           "检查规则入参错误}}",
	"_valueTypeError":       "参数类型不正确}}",
	"_fieldNotFound":        "{{字段不存在}}",
}

// NewErrorContext return a error context.
func NewErrorContext(c CheckerContext) *ErrContext {
	errCtx := &ErrContext{
		FieldName:  c.FieldName,
		FieldLabel: c.FieldLabel,
		FieldValue: c.FieldValue,

		errorTemplate: getErrorTemplate(c.Rule.checker),
	}
	errCtx.makeMessage()

	return errCtx
}

func (e *ErrContext) makeMessage() {
	msg := e.errorTemplate
	if strings.Contains(e.errorTemplate, FieldNamePlaceholder) || strings.Contains(e.errorTemplate, FieldLimitPlaceholder) {
		msg = strings.NewReplacer(
			FieldNamePlaceholder, e.FieldName,
			FieldLimitPlaceholder, fmt.Sprintf("%v", e.fieldLimitValue),
		).Replace(msg)
	}

	fieldLabelPrefix, limitValueSuffix := !strings.HasPrefix(e.errorTemplate, "{{"), !strings.HasSuffix(e.errorTemplate, "}}")
	if fieldLabelPrefix {
		msg = e.FieldLabel + msg
	} else {
		msg = msg[2:] // Remove the first two "{{"
	}
	if limitValueSuffix {
		if e.fieldLimitValue != nil {
			msg += fmt.Sprintf("%v", e.fieldLimitValue)
		}
	} else {
		msg = msg[:len(msg)-2] // Remove the last "}}"
	}

	e.errorMessage = msg
}

func (e *ErrContext) SetFieldLimitValue(v interface{}) {
	e.fieldLimitValue = v
	e.makeMessage()
}

func (e *ErrContext) SetTemplate(key string) {
	e.errorTemplate = getErrorTemplate(key)
	e.makeMessage()
}

// getErrorTemplate return the template of the given rule name.
func getErrorTemplate(key string) string {
	if value, ok := errorTemplate[key]; ok {
		return value
	}
	return errorTemplate["_unknownErrorTemplate"]
}

// SetMessageTemplates set message templates.
func SetMessageTemplates(messages map[string]string) {
	for k, v := range messages {
		errorTemplate[k] = v
	}
}

func MakeUserDefinedError(msg string) *ErrContext {
	errCtx := &ErrContext{
		errorMessage: msg,
	}
	return errCtx
}

func MakeCheckerNotFoundError(c CheckerContext) *ErrContext {
	template := strings.TrimPrefix(getErrorTemplate("_checkerNotFound"), "~")

	errCtx := &ErrContext{
		FieldName:       c.FieldName,
		FieldLabel:      c.FieldLabel,
		FieldValue:      c.FieldValue,
		fieldLimitValue: c.Rule.params,
		errorTemplate:   template,
	}
	errCtx.makeMessage()
	return errCtx
}

func MakeCheckerParamError(c CheckerContext) *ErrContext {
	template := strings.TrimPrefix(getErrorTemplate("_paramError"), "~")

	errCtx := &ErrContext{
		FieldName:       c.FieldName,
		FieldLabel:      c.FieldLabel,
		FieldValue:      c.FieldValue,
		fieldLimitValue: c.Rule.params,
		errorTemplate:   template,
	}
	errCtx.makeMessage()
	return errCtx
}

func MakeValueTypeError(c CheckerContext) *ErrContext {
	template := strings.TrimPrefix(getErrorTemplate("_valueTypeError"), "~")

	errCtx := &ErrContext{
		FieldName:       c.FieldName,
		FieldLabel:      c.FieldLabel,
		FieldValue:      c.FieldValue,
		fieldLimitValue: c.Rule.params,
		errorTemplate:   template,
	}
	errCtx.makeMessage()
	return errCtx
}

func MakeFieldNotFoundError(c CheckerContext) *ErrContext {
	template := strings.TrimPrefix(getErrorTemplate("_fieldNotFound"), "~")

	errCtx := &ErrContext{
		FieldName:       c.FieldName,
		FieldLabel:      c.FieldLabel,
		FieldValue:      c.FieldValue,
		fieldLimitValue: c.Rule.params,
		errorTemplate:   template,
	}
	errCtx.makeMessage()
	return errCtx
}
