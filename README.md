# govalid
[![Build Status](https://travis-ci.com/wuhan005/govalid.svg?branch=master)](https://travis-ci.com/wuhan005/govalid)
[![Go Report Card](https://goreportcard.com/badge/github.com/wuhan005/govalid)](https://goreportcard.com/report/github.com/wuhan005/govalid)

## How to use
1. `go get github.com/wuhan005/govalid`
2. 
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

## LICENSE
MIT