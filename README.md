# govalid


## How to use
1. `go get github.com/wuhan005/govalid`
2. 
```go
r := registerForm struct {
	Name string `valid:"required;username" label:"昵称"`
	ID   int    `valid:"required;min:0;max:999" label:"用户编号"`
	Mail string `valid:"email" label:"电子邮箱"`
	Card string `valid:"idcard" label:""`
}

v := govalid.New(r)
if !v.Check(){
    for _, v := range v.errors {
        fmt.Println(v.Message)
    }
}
```

## LICENSE
MIT