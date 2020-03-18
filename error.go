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
	"_ruleNotFound":   "检查规则未找到",
	"_unknown":        "未知错误",
	"_paramError":     "检查规则入参错误",
	"_valueTypeError": "参数类型不正确",
}

func getErrorTemplate(key string) string {
	if value, ok := errorTemplate[key]; ok {
		return value
	} else {
		return errorTemplate["_unknown"]
	}
}
