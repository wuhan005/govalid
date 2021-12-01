# govalid

A simple Go struct validator.

[![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/govalid)](https://goreportcard.com/report/github.com/wuhan005/govalid)

## Quick Start

`go get -u github.com/wuhan005/govalid`

```
v := struct {
    Name string `valid:"required;username" label:"昵称"`
    ID   int    `valid:"required;min:0;max:999" label:"用户编号"`
    Mail string `valid:"required;email" label:""`
}{
    "e99_", 1990, "i@github.red.",
}

errs, ok := govalid.Check(v)
if !ok {
    for _, err := range errs {
        fmt.Println(err)
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

```
govalid.SetMessageTemplates(map[string]string{
    "required": "can't be null.",
    "min": "must bigger than %v.",
})
```

## Customize Error Check Function

```
func main() {
	// set your error message.
	govalid.SetMessageTemplates(map[string]string{
		"mycheck": "content cann't contain 'e99'.",
	})

	// add new check function
	govalid.Checkers["mycheck"] = func(c govalid.CheckerContext) *govalid.ErrContext {
		errCtx := govalid.NewErrorContext(c)
		value, ok := ctx.Value.(string)
		if !ok {
			return MakeCheckerParamError(c)
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

    errs, ok := govalid.Check(v)
    if !ok {
        for _, err := range errs {
            fmt.Println(err)
        }
    }
}
```

## Special Thanks

https://github.com/jiazhoulvke/echo-form

## LICENSE

MIT
