# govalid

A simple Go struct validator.

[![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/govalid)](https://goreportcard.com/report/github.com/wuhan005/govalid)

## Quick Start

```bash
go get -u github.com/wuhan005/govalid
```

```go
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

```text
昵称的最后一个字符不能为下划线
用户编号应小于999
不是合法的电子邮箱格式
```

## Customize Error Message

`SetMessageTemplates` merges the given templates into the current language's
template set, overriding any existing entries.

```go
govalid.SetMessageTemplates(map[string]string{
    "required": "can not be null",
    "min":      "must be greater than",
})
```

You can also pass a language tag to update a specific locale:

```go
govalid.SetMessageTemplates(map[string]string{
    "required": "must not be empty",
}, language.English)
```

## Customize Error Check Function

```go
package main

import (
    "fmt"
    "strings"

    "github.com/wuhan005/govalid"
)

func main() {
    // Set the error message for the new check rule.
    govalid.SetMessageTemplates(map[string]string{
        "mycheck": "content can not contain 'e99'",
    })

    // Add a new check function.
    govalid.Checkers["mycheck"] = func(c govalid.CheckerContext) *govalid.ErrContext {
        value, ok := c.FieldValue.(string)
        if !ok {
            return govalid.MakeValueTypeError(c)
        }
        if strings.Contains(value, "e99") {
            return govalid.NewErrorContext(c)
        }
        // If the check passed, return nil.
        return nil
    }

    r := struct {
        Content string `valid:"mycheck"`
    }{
        Content: "helloe99",
    }

    errs, ok := govalid.Check(r)
    if !ok {
        for _, err := range errs {
            fmt.Println(err)
        }
    }
}
```

## Custom `Validate` Method

Any struct that implements a `Validate() error` method is automatically called
after the tag-based rules. Use it for cross-field or business rule validation:

```go
type Form struct {
    Password       string `valid:"required" label:"密码"`
    RepeatPassword string `valid:"required" label:"重复密码"`
}

func (f *Form) Validate() error {
    if f.Password != f.RepeatPassword {
        return errors.New("两次输入的密码不一致")
    }
    return nil
}
```

## Special Thanks

https://github.com/jiazhoulvke/echo-form

## LICENSE

MIT
