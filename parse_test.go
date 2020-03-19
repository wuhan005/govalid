package govalid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type registerForm struct {
	Name string `valid:"required;username" label:"昵称"`
	ID   int    `valid:"required;min:5;max:10" label:"用户编号"`
	Mail string `valid:"email" label:"电子邮箱"`
	Card string `valid:"idcard" label:""`
}

func TestNew(t *testing.T) {
	s1 := struct {
		Name string `valid:`
	}{}
	s2 := struct {
		Name string `valid:""`
	}{}
	s3 := struct {
		Name string `valid:"abc"`
	}{}
	s4 := struct {
		Name string `valid111:"abc"`
	}{}

	assert.Equal(t, New(s1).Check(), true)
	assert.Equal(t, New(s2).Check(), true)
	assert.Equal(t, New(s3).Check(), false)
	assert.Equal(t, New(s4).Check(), true)
}
