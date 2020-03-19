package govalid

import (
	"fmt"
	"testing"
)

type registerForm struct {
	Name string `valid:"required;username" label:"昵称"`
	ID   int    `valid:"required;min:5;max:10" label:"用户编号"`
	Mail string `valid:"email" label:"电子邮箱"`
	Card string `valid:"idcard" label:""`
}

func TestNew(t *testing.T) {
	r := registerForm{
		Name: "e99e99",
		ID:   9,
		Mail: "e99@e.99",
		Card: "1231232",
	}
	v := New(r)
	if !v.Check() {
		for _, err := range v.errors {
			fmt.Println(err.Message)
		}
	}
}
