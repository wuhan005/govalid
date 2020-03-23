# govalid
[![Build Status](https://travis-ci.com/wuhan005/govalid.svg?branch=master)](https://travis-ci.com/wuhan005/govalid)
[![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/govalid)](https://goreportcard.com/report/github.com/wuhan005/govalid)

## Quick Start
`go get -u github.com/wuhan005/govalid`

```go
r := struct {
    Name string `valid:"required;username" label:"昵称"`
    ID   int    `valid:"required;min:0;max:999" label:"用户编号"`
    Mail string `valid:"required;email" label:""`
}{
    "e99_", 1990, "i@github.red.",
}

v := govalid.New(r)
if !v.Check() {
    for _, err := range v.Errors {
        fmt.Println(err.Message)
    }
}
```
Output:
```
昵称的最后一个字符不能为下划线
用户编号应小于999
不是合法的电子邮箱格式
```

## Customize Error Message
```go
govalid.SetDefaultMessage(map[string]string{
    "required": "can't be null.",
    "min": "must bigger than %v.",
})
```

### Default Error Messages
```
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
```
## Customize Error Check Function
```go
func main() {
	// set your error message.
	govalid.SetDefaultMessage(map[string]string{
		"mycheck": "content cann't contain 'e99'.",
	})

	// add new check function
	govalid.Checkers["mycheck"] = func(c govalid.RuleContext) *govalid.ErrContext {
		ctx := govalid.NewErrorContext(c)
		value, ok := ctx.Value.(string)
		if !ok {
			ctx.SetMessage("Input must be string!")
			return ctx
		}
		if strings.Contains(value, "e99") {
			return ctx
		}
		// if the check passed, return nil.
		return nil
	}

	r := struct {
		Content string `valid:"mycheck"`
	}{
		"helloe99",
	}

	v := govalid.New(r)
	if !v.Check() {
		for _, err := range v.Errors {
			fmt.Println(err.Message)
		}
	}
}
```
## Special Thanks
https://github.com/jiazhoulvke/echo-form

## LICENSE
MIT