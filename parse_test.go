package govalid

import (
	"fmt"
	"testing"
)

type registerForm struct {
	Name string `valid:"required;alpha" label:"昵称"`
	ID   int    `valid:"required;min:5;max:10" label:"用户编号"`
}

func TestNew(t *testing.T) {
	r := registerForm{
		Name: "e99",
		ID:   11,
	}
	v := New(r)
	if !v.Check() {
		for _, err := range v.errors {
			fmt.Println(err.Label + err.Message)
		}
	}
}
