# govalid

A simple, struct-tag driven validator for Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/wuhan005/govalid.svg)](https://pkg.go.dev/github.com/wuhan005/govalid)
[![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/govalid)](https://goreportcard.com/report/github.com/wuhan005/govalid)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

- **Tag-based** — declare rules right next to your struct fields.
- **Composable** — chain multiple checkers per field with `;`.
- **i18n ready** — Chinese & English bundled, any locale pluggable.
- **Extensible** — register custom checkers, override messages, hook a `Validate()` method for cross-field rules.
- **Safe by default** — handles nil pointers, maps, slices, embedded structs and unexported fields without panicking.
- **Battle-tested** — 96%+ test coverage and a fuzz suite.

## Install

```bash
go get -u github.com/wuhan005/govalid
```

Requires Go 1.16+.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/wuhan005/govalid"
)

func main() {
    v := struct {
        Name string `valid:"required;username" label:"昵称"`
        ID   int    `valid:"required;min:0;max:999" label:"用户编号"`
        Mail string `valid:"required;email" label:"邮箱"`
    }{
        Name: "e99_",
        ID:   1990,
        Mail: "i@github.red.",
    }

    errs, ok := govalid.Check(v)
    if !ok {
        for _, err := range errs {
            fmt.Println(err)
        }
    }
}
```

```text
昵称的最后一个字符不能为下划线
用户编号应小于999
邮箱不是合法的电子邮箱格式
```

## Tags

govalid reads three struct tags:

| Tag | Default | Purpose |
| --- | --- | --- |
| `valid` | — | Validation rules, separated by `;`. Parameters follow `:`, multiple parameters separated by `,`. |
| `label` | field name | Human-friendly name used in error messages. Supports per-locale overrides via `label-en`, `label-zh`, … |
| `msg` | — | Override the entire error message for this field. The first failing rule short-circuits to this message. |

The tag names themselves are configurable through the package-level
`govalid.RulesField`, `govalid.LabelField` and `govalid.MessageField`
variables if you need to coexist with another tag namespace.

```go
type Form struct {
    Email string `valid:"required;email" label:"邮箱" label-en:"Email" msg:"邮箱格式不对"`
}
```

## Built-in Checkers

| Rule | Parameters | Applies to | Description |
| --- | --- | --- | --- |
| `required` | — | any | Field must be non-zero. Slices, arrays, maps, strings and channels must be non-empty; pointers, interfaces and funcs must be non-nil. |
| `min:N` | one number | int / uint / float | Numeric lower bound (inclusive). |
| `max:N` | one number | int / uint / float | Numeric upper bound (inclusive). |
| `minlen:N` | one int | string / slice / array / map | Minimum length. Strings are counted in **runes**, not bytes. |
| `maxlen:N` | one int | string / slice / array / map | Maximum length. |
| `alpha` | — | string | ASCII letters only (`a-z`, `A-Z`). |
| `alphanumeric` | — | string | ASCII letters or digits. |
| `alphadash` | — | string | Letters, digits, or underscores (`\w+`). |
| `username` | — | string | Combination of `alphadash` plus first-char-must-be-letter and no trailing underscore. |
| `email` | — | string | Email address syntax. |
| `ipv4` | — | string | IPv4 address syntax. |
| `mobile` | — | string | Chinese mobile phone number (with optional `+86` / `86` prefix). |
| `tel` | — | string | Chinese landline number. |
| `phone` | — | string | Either `mobile` or `tel`. |
| `idcard` | — | string | Chinese 15- or 18-digit ID card number (final character `0-9`, `X`, or `x`). |
| `equal:OtherField` | one field name | any | Stringified value must match another sibling field. |
| `list:a,b,c` | one or more values | any | Stringified value must be one of the listed values. |

Empty strings short-circuit to "ok" for all string-format checkers
(`alpha`, `email`, `ipv4`, `mobile`, …) so you can opt fields in and out
by combining them with `required`:

```go
type Form struct {
    Optional string `valid:"email"`            // empty allowed
    Required string `valid:"required;email"`   // empty rejected
}
```

## Cross-field & Business Rules — `Validate() error`

Anything more complex than a single field belongs in a `Validate()`
method. govalid invokes it after the tag rules, on both value and
pointer receivers:

```go
type Form struct {
    Password       string `valid:"required;minlen:8" label:"密码"`
    RepeatPassword string `valid:"required" label:"重复密码"`
}

func (f *Form) Validate() error {
    if f.Password != f.RepeatPassword {
        return errors.New("两次输入的密码不一致")
    }
    return nil
}
```

The signature must be exactly `func() error`; anything else is silently
ignored.

## Nested Structs & Slices

Nested structs and slices of structs are walked automatically — every
field's tags fire just like top-level fields:

```go
type Item struct {
    Name string `valid:"required" label:"项目名"`
    Qty  int    `valid:"min:1" label:"数量"`
}

type Cart struct {
    Items []Item `valid:"required" label:"购物车"`
}
```

Embedded (anonymous) structs are also fully supported.

## Customizing Error Messages

`SetMessageTemplates` merges your templates into a locale's template
set. Without a locale argument it updates the default locale (Chinese).

```go
govalid.SetMessageTemplates(map[string]string{
    "required": "can not be null",
    "min":      "must be greater than",
})

govalid.SetMessageTemplates(map[string]string{
    "required": "must not be empty",
}, language.English)
```

Per-call locale selection happens through `Check`'s variadic argument:

```go
errs, ok := govalid.Check(form, language.English)
```

Unknown locales fall back to the default (Chinese) template set; missing
keys fall back to a generic "unknown error" template.

## Adding Your Own Checker

Register a function in the `govalid.Checkers` map. The error helpers
(`NewErrorContext`, `MakeValueTypeError`, `MakeCheckerParamError`,
`MakeFieldNotFoundError`) take care of formatting:

```go
package main

import (
    "fmt"
    "strings"

    "github.com/wuhan005/govalid"
)

func main() {
    govalid.SetMessageTemplates(map[string]string{
        "noE99": "can not contain 'e99'",
    })

    govalid.Checkers["noE99"] = func(c govalid.CheckerContext) *govalid.ErrContext {
        v, ok := c.FieldValue.(string)
        if !ok {
            return govalid.MakeValueTypeError(c)
        }
        if strings.Contains(v, "e99") {
            return govalid.NewErrorContext(c)
        }
        return nil
    }

    r := struct {
        Content string `valid:"noE99" label:"内容"`
    }{Content: "helloe99"}

    if errs, ok := govalid.Check(r); !ok {
        for _, err := range errs {
            fmt.Println(err)
        }
    }
}
```

## API Reference

```go
// Check validates v and returns the list of failures together with a
// boolean indicating overall success. Pass an optional language tag to
// pick a non-default locale.
func Check(v interface{}, lang ...language.Tag) (errs []*ErrContext, ok bool)

// SetMessageTemplates merges templates into the given locale (default
// when omitted), overriding existing entries.
func SetMessageTemplates(templates map[string]string, lang ...language.Tag)

// Checkers is the registry of validation functions, keyed by rule name.
var Checkers map[string]CheckFunc

// Tag names — change these once at startup if you need different keys.
var RulesField, LabelField, MessageField string
```

`*ErrContext` implements the `error` interface, so each item can be
returned, wrapped, or logged like any other Go error.

## Special Thanks

[github.com/jiazhoulvke/echo-form](https://github.com/jiazhoulvke/echo-form)

## License

[MIT](LICENSE)
