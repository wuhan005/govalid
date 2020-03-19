package govalid

type errContext struct {
	Field      string // form field name
	Label      string
	Tmpl       string // error message template
	Message    string
	Value      interface{}
	LimitValue interface{}
}

var errorTemplate = map[string]string{
	"required":        "不能为空",
	"min":             "应大于%v",
	"max":             "应小于%v",
	"alpha":           "必须只包含字母",
	"alphanumeric":    "只能含有字母或数字",
	"alphaDash":       "只含有数字或字母以及下划线",
	"firstCharAlpha":  "的第一个字符必须为字母",
	"lastUnderline":   "最后一个字符不能为下划线",
	"_ruleNotFound":   "检查规则未找到",
	"_unknown":        "未知错误",
	"_paramError":     "检查规则入参错误",
	"_valueTypeError": "参数类型不正确",
}

func NewErrorContext(c ruleContext) *errContext {
	return &errContext{
		Tmpl:       getErrorTemplate(c.rule),
		Message:    getErrorTemplate(c.rule), // default message is the raw template
		Value:      c.value,
		LimitValue: c.params,
	}
}

func (e *errContext) SetMessage(msg string) {
	e.Message = msg
}

func getErrorTemplate(key string) string {
	if value, ok := errorTemplate[key]; ok {
		return value
	} else {
		return errorTemplate["_unknown"]
	}
}
