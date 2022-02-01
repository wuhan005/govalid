package govalid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_required(t *testing.T) {
	v := struct {
		Name string `valid:"required" label:"用户名"`
	}{
		"",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "用户名不能为空", errs[0].Error())

	v = struct {
		Name string `valid:"required" label:"用户名"`
	}{
		"E99p1ant",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, errs)
}

func Test_min(t *testing.T) {
	v := struct {
		Score int `valid:"min:0" label:"评分"`
	}{
		-233,
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "评分应大于0", errs[0].Error())

	v = struct {
		Score int `valid:"min:0" label:"评分"`
	}{
		233,
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_max(t *testing.T) {
	v := struct {
		Score int `valid:"max:100" label:"得分"`
	}{
		1024,
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "得分应小于100", errs[0].Error())

	v = struct {
		Score int `valid:"max:100" label:"得分"`
	}{
		47,
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_minlen(t *testing.T) {
	v := struct {
		Message string `valid:"minlen:5" label:"留言"`
	}{
		"aaa",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "留言长度应大于5", errs[0].Error())

	v1 := struct {
		Message string `valid:"minlen:5.2" label:"留言"`
	}{
		"aaa",
	}
	errs, ok = Check(v1)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "留言检查规则入参错误", errs[0].Error())

	v3 := struct {
		Message string `valid:"minlen:5" label:"留言"`
	}{
		"",
	}
	errs, ok = Check(v3)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v4 := struct {
		Message string `valid:"minlen:5" label:"留言"`
	}{
		"Hello e99!",
	}
	errs, ok = Check(v4)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_maxlen(t *testing.T) {
	v := struct {
		Message string `valid:"maxlen:8" label:"留言"`
	}{
		"this_is_e99999",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "留言长度应小于8", errs[0].Error())

	v1 := struct {
		Message string `valid:"maxlen:5.2" label:"留言"`
	}{
		"aaa",
	}
	errs, ok = Check(v1)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "留言检查规则入参错误", errs[0].Error())

	v2 := struct {
		Message string `valid:"maxlen:8" label:"留言"`
	}{
		"",
	}
	errs, ok = Check(v2)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v3 := struct {
		Message string `valid:"maxlen:8" label:"留言"`
	}{
		"e99",
	}
	errs, ok = Check(v3)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_alpha(t *testing.T) {
	v := struct {
		Name string `valid:"alpha" label:"昵称"`
	}{
		"e99p1ant",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "昵称必须只包含字母", errs[0].Error())

	v = struct {
		Name string `valid:"alpha" label:"昵称"`
	}{
		"eggplant",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v = struct {
		Name string `valid:"alpha" label:"昵称"`
	}{
		"",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_alphanumeric(t *testing.T) {
	v := struct {
		Name string `valid:"alphanumeric" label:"昵称"`
	}{
		"e99p|ant",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "昵称只能含有字母或数字", errs[0].Error())

	v = struct {
		Name string `valid:"alphanumeric" label:"昵称"`
	}{
		"e99p1ant",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v = struct {
		Name string `valid:"alphanumeric" label:"昵称"`
	}{
		"",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_alphaDash(t *testing.T) {
	v := struct {
		Name string `valid:"alphadash" label:"昵称"`
	}{
		"e99p_1ant$",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "昵称只含有数字或字母以及下划线", errs[0].Error())

	v = struct {
		Name string `valid:"alphadash" label:"昵称"`
	}{
		"e99p1__ant__",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v = struct {
		Name string `valid:"alphadash" label:"昵称"`
	}{
		"",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_userName(t *testing.T) {
	v := struct {
		Name string `valid:"username" label:"昵称"`
	}{
		"199p1ant",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "昵称的第一个字符必须为字母", errs[0].Error())

	v = struct {
		Name string `valid:"username" label:"昵称"`
	}{
		"e99p1ant_",
	}
	errs, ok = Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "昵称的最后一个字符不能为下划线", errs[0].Error())

	v = struct {
		Name string `valid:"username" label:"昵称"`
	}{
		"e99p1ant",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v = struct {
		Name string `valid:"username" label:"昵称"`
	}{
		"",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_email(t *testing.T) {
	v := struct {
		Email string `valid:"email" label:""`
	}{
		"e99@q.",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "不是合法的电子邮箱格式", errs[0].Error())

	v1 := struct {
		Email string `valid:"email" label:"Mailll"`
	}{
		"e99@@@@99.com",
	}
	errs, ok = Check(v1)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "Mailll不是合法的电子邮箱格式", errs[0].Error())

	v = struct {
		Email string `valid:"email" label:""`
	}{
		"e99@q.a.a.a.a.a.aa.a.a.com",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v = struct {
		Email string `valid:"email" label:""`
	}{
		"",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_ipv4(t *testing.T) {
	v := struct {
		IP string `valid:"ipv4" label:""`
	}{
		"1.2.3.256",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "不是合法的 IPv4 地址格式", errs[0].Error())

	v2 := struct {
		IP string `valid:"ipv4" label:"IPIPIPP"`
	}{
		"255.255.255.255.",
	}
	errs, ok = Check(v2)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "IPIPIPP不是合法的 IPv4 地址格式", errs[0].Error())

	v = struct {
		IP string `valid:"ipv4" label:""`
	}{
		"255.255.255.255",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v = struct {
		IP string `valid:"ipv4" label:""`
	}{
		"127.128.129.130",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v = struct {
		IP string `valid:"ipv4" label:""`
	}{
		"",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_mobile(t *testing.T) {
	v := struct {
		Phone string `valid:"mobile" label:"电话号码"`
	}{
		"13888888888a",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "电话号码不是合法的手机号", errs[0].Error())

	v = struct {
		Phone string `valid:"mobile" label:"电话号码"`
	}{
		"1388888888",
	}
	errs, ok = Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "电话号码不是合法的手机号", errs[0].Error())

	v = struct {
		Phone string `valid:"mobile" label:"电话号码"`
	}{
		"1088888888",
	}
	errs, ok = Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "电话号码不是合法的手机号", errs[0].Error())

	v = struct {
		Phone string `valid:"mobile" label:"电话号码"`
	}{
		"2388888888",
	}
	errs, ok = Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "电话号码不是合法的手机号", errs[0].Error())

	v = struct {
		Phone string `valid:"mobile" label:"电话号码"`
	}{
		"13888888888",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_tel(t *testing.T) {
	v := struct {
		Phone string `valid:"tel" label:"电话号码"`
	}{
		"13888888888",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "电话号码不是合法的座机号码", errs[0].Error())

	v = struct {
		Phone string `valid:"tel" label:"电话号码"`
	}{
		"qqqqq",
	}
	errs, ok = Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "电话号码不是合法的座机号码", errs[0].Error())

	v = struct {
		Phone string `valid:"tel" label:"电话号码"`
	}{
		"1111111a",
	}
	errs, ok = Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "电话号码不是合法的座机号码", errs[0].Error())

	v = struct {
		Phone string `valid:"tel" label:"电话号码"`
	}{
		"26088888",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v = struct {
		Phone string `valid:"tel" label:"电话号码"`
	}{
		"47474747",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_phone(t *testing.T) {
	v := struct {
		Phone string `valid:"phone" label:"电话号码"`
	}{
		"@#$%^&*",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "电话号码不是合法的号码", errs[0].Error())

	v = struct {
		Phone string `valid:"phone" label:"电话号码"`
	}{
		"123456",
	}
	errs, ok = Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "电话号码不是合法的号码", errs[0].Error())

	v = struct {
		Phone string `valid:"phone" label:"电话号码"`
	}{
		"1111111a",
	}
	errs, ok = Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "电话号码不是合法的号码", errs[0].Error())

	v = struct {
		Phone string `valid:"phone" label:"电话号码"`
	}{
		"13888888888",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v = struct {
		Phone string `valid:"phone" label:"电话号码"`
	}{
		"47474747",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_idCard(t *testing.T) {
	v := struct {
		Phone string `valid:"idcard" label:"身份证号"`
	}{
		"12312312312312",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "身份证号不是合法的身份证号", errs[0].Error())

	v = struct {
		Phone string `valid:"idcard" label:"身份证号"`
	}{
		"@#$%^&*",
	}
	errs, ok = Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "身份证号不是合法的身份证号", errs[0].Error())

	v = struct {
		Phone string `valid:"idcard" label:"身份证号"`
	}{
		"1234567890987654XX",
	}
	errs, ok = Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "身份证号不是合法的身份证号", errs[0].Error())

	v = struct {
		Phone string `valid:"idcard" label:"身份证号"`
	}{
		"123456789098765432",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))

	v = struct {
		Phone string `valid:"idcard" label:"身份证号"`
	}{
		"12345678909876543X",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
	assert.Zero(t, len(errs))
}

func Test_Equal(t *testing.T) {
	v := struct {
		Password       string `valid:"required" label:"密码"`
		RepeatPassword string `valid:"equal:Password" label:"重复密码"`
	}{
		Password:       "123456",
		RepeatPassword: "1234567",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "重复密码的值前后不相同", errs[0].Error())

	v = struct {
		Password       string `valid:"required" label:"密码"`
		RepeatPassword string `valid:"equal:Password" label:"重复密码"`
	}{
		Password:       "123456",
		RepeatPassword: "123456",
	}
	errs, ok = Check(v)
	assert.True(t, ok)
}
