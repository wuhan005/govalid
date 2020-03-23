package govalid

// ErrContext contains the error context.
type ErrContext struct {
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
	"minlen":          "长度应大于%v",
	"maxlen":          "长度应小于%v",
	"alpha":           "必须只包含字母",
	"alphanumeric":    "只能含有字母或数字",
	"alphadash":       "只含有数字或字母以及下划线",
	"firstCharAlpha":  "的第一个字符必须为字母",
	"lastUnderline":   "的最后一个字符不能为下划线",
	"email":           "不是合法的电子邮箱格式",
	"ipv4":            "不是合法的 IPv4 地址格式",
	"mobile":          "不是合法的手机号",
	"tel":             "不是合法的座机号码",
	"phone":           "不是合法的号码",
	"idcard":          "不是合法的身份证号",
	"_ruleNotFound":   "检查规则未找到",
	"_unknown":        "未知错误",
	"_paramError":     "检查规则入参错误",
	"_valueTypeError": "参数类型不正确",
}

// NewErrorContext return a error context.
func NewErrorContext(c RuleContext) *ErrContext {
	return &ErrContext{
		Tmpl:       GetErrorTemplate(c.rule),
		Message:    GetErrorTemplate(c.rule), // default message is the raw template
		Value:      c.value,
		LimitValue: c.params,
	}
}

// SetMessage set the error context's message.
func (e *ErrContext) SetMessage(msg string) {
	e.Message = msg
}

// GetErrorTemplate return the template of giving rule name.
func GetErrorTemplate(key string) string {
	if value, ok := errorTemplate[key]; ok {
		return value
	}
	return errorTemplate["_unknown"]
}

// SetDefaultMessage set the default message template.
func SetDefaultMessage(msg map[string]string) {
	if len(msg) == 0 {
		return
	}

	for name := range msg {
		errorTemplate[name] = msg[name]
	}
}
