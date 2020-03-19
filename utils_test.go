package govalid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_required(t *testing.T) {
	s := struct {
		Name string `valid:"required" label:"用户名"`
	}{
		"",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "用户名不能为空")

	s = struct {
		Name string `valid:"required" label:"用户名"`
	}{
		"e99",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_min(t *testing.T) {
	s := struct {
		Score int `valid:"min:0" label:"评分"`
	}{
		-233,
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "评分应大于0")

	s = struct {
		Score int `valid:"min:0" label:"评分"`
	}{
		233,
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_max(t *testing.T) {
	s := struct {
		Score int `valid:"max:100" label:"得分"`
	}{
		1024,
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "得分应小于100")

	s = struct {
		Score int `valid:"max:100" label:"得分"`
	}{
		47,
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_minlen(t *testing.T) {
	s := struct {
		Message string `valid:"minlen:5" label:"留言"`
	}{
		"aaa",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "留言长度应大于5")

	s1 := struct {
		Message string `valid:"minlen:5.2" label:"留言"`
	}{
		"aaa",
	}
	v = New(s1)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "留言检查规则入参错误")

	s = struct {
		Message string `valid:"minlen:5" label:"留言"`
	}{
		"Hello e99!",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_maxlen(t *testing.T) {
	s := struct {
		Message string `valid:"maxlen:8" label:"留言"`
	}{
		"this_is_e99999",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "留言长度应小于8")

	s1 := struct {
		Message string `valid:"maxlen:5.2" label:"留言"`
	}{
		"aaa",
	}
	v = New(s1)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "留言检查规则入参错误")

	s = struct {
		Message string `valid:"maxlen:8" label:"留言"`
	}{
		"e99",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_alpha(t *testing.T) {
	s := struct {
		Name string `valid:"alpha" label:"昵称"`
	}{
		"e99p1ant",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "昵称必须只包含字母")

	s = struct {
		Name string `valid:"alpha" label:"昵称"`
	}{
		"eggplant",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

	s = struct {
		Name string `valid:"alpha" label:"昵称"`
	}{
		"",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_alphanumeric(t *testing.T) {
	s := struct {
		Name string `valid:"alphanumeric" label:"昵称"`
	}{
		"e99p|ant",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "昵称只能含有字母或数字")

	s = struct {
		Name string `valid:"alphanumeric" label:"昵称"`
	}{
		"e99p1ant",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

	s = struct {
		Name string `valid:"alphanumeric" label:"昵称"`
	}{
		"",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

}

func Test_alphaDash(t *testing.T) {
	s := struct {
		Name string `valid:"alphadash" label:"昵称"`
	}{
		"e99p_1ant$",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "昵称只含有数字或字母以及下划线")

	s = struct {
		Name string `valid:"alphadash" label:"昵称"`
	}{
		"e99p1__ant__",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

	s = struct {
		Name string `valid:"alphadash" label:"昵称"`
	}{
		"",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_userName(t *testing.T) {
	s := struct {
		Name string `valid:"username" label:"昵称"`
	}{
		"199p1ant",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "昵称的第一个字符必须为字母")

	s = struct {
		Name string `valid:"username" label:"昵称"`
	}{
		"e99p1ant_",
	}
	v = New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "昵称的最后一个字符不能为下划线")

	s = struct {
		Name string `valid:"username" label:"昵称"`
	}{
		"e99p1ant",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

	s = struct {
		Name string `valid:"username" label:"昵称"`
	}{
		"",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_email(t *testing.T) {
	s := struct {
		Email string `valid:"email" label:""`
	}{
		"e99@q.",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "不是合法的电子邮箱格式")

	s2 := struct {
		Email string `valid:"email" label:"Mailll"`
	}{
		"e99@@@@99.com",
	}
	v = New(s2)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "Mailll不是合法的电子邮箱格式")

	s = struct {
		Email string `valid:"email" label:""`
	}{
		"e99@q.a.a.a.a.a.aa.a.a.com",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

	s = struct {
		Email string `valid:"email" label:""`
	}{
		"",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_ipv4(t *testing.T) {
	s := struct {
		IP string `valid:"ipv4" label:""`
	}{
		"1.2.3.256",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "不是合法的 IPv4 地址格式")

	s2 := struct {
		IP string `valid:"ipv4" label:"IPIPIPP"`
	}{
		"255.255.255.255.",
	}
	v = New(s2)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "IPIPIPP不是合法的 IPv4 地址格式")

	s = struct {
		IP string `valid:"ipv4" label:""`
	}{
		"255.255.255.255",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

	s = struct {
		IP string `valid:"ipv4" label:""`
	}{
		"127.128.129.130",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

	s = struct {
		IP string `valid:"ipv4" label:""`
	}{
		"",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_mobile(t *testing.T) {
	s := struct {
		Phone string `valid:"mobile" label:"电话号码"`
	}{
		"13888888888a",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "电话号码不是合法的手机号")

	s = struct {
		Phone string `valid:"mobile" label:"电话号码"`
	}{
		"1388888888",
	}
	v = New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "电话号码不是合法的手机号")

	s = struct {
		Phone string `valid:"mobile" label:"电话号码"`
	}{
		"1088888888",
	}
	v = New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "电话号码不是合法的手机号")

	s = struct {
		Phone string `valid:"mobile" label:"电话号码"`
	}{
		"2388888888",
	}
	v = New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "电话号码不是合法的手机号")

	s = struct {
		Phone string `valid:"mobile" label:"电话号码"`
	}{
		"13888888888",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_tel(t *testing.T) {
	s := struct {
		Phone string `valid:"tel" label:"电话号码"`
	}{
		"13888888888",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "电话号码不是合法的座机号码")

	s = struct {
		Phone string `valid:"tel" label:"电话号码"`
	}{
		"qqqqq",
	}
	v = New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "电话号码不是合法的座机号码")

	s = struct {
		Phone string `valid:"tel" label:"电话号码"`
	}{
		"1111111a",
	}
	v = New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "电话号码不是合法的座机号码")

	s = struct {
		Phone string `valid:"tel" label:"电话号码"`
	}{
		"26088888",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

	s = struct {
		Phone string `valid:"tel" label:"电话号码"`
	}{
		"47474747",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_phone(t *testing.T) {
	s := struct {
		Phone string `valid:"phone" label:"电话号码"`
	}{
		"@#$%^&*",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "电话号码不是合法的号码")

	s = struct {
		Phone string `valid:"phone" label:"电话号码"`
	}{
		"123456",
	}
	v = New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "电话号码不是合法的号码")

	s = struct {
		Phone string `valid:"phone" label:"电话号码"`
	}{
		"1111111a",
	}
	v = New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "电话号码不是合法的号码")

	s = struct {
		Phone string `valid:"phone" label:"电话号码"`
	}{
		"13888888888",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

	s = struct {
		Phone string `valid:"phone" label:"电话号码"`
	}{
		"47474747",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}

func Test_idCard(t *testing.T) {
	s := struct {
		Phone string `valid:"idcard" label:"身份证号"`
	}{
		"12312312312312",
	}
	v := New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "身份证号不是合法的身份证号")

	s = struct {
		Phone string `valid:"idcard" label:"身份证号"`
	}{
		"@#$%^&*",
	}
	v = New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "身份证号不是合法的身份证号")

	s = struct {
		Phone string `valid:"idcard" label:"身份证号"`
	}{
		"1234567890987654XX",
	}
	v = New(s)
	assert.Equal(t, v.Check(), false)
	assert.Equal(t, v.Errors[0].Message, "身份证号不是合法的身份证号")

	s = struct {
		Phone string `valid:"idcard" label:"身份证号"`
	}{
		"123456789098765432",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)

	s = struct {
		Phone string `valid:"idcard" label:"身份证号"`
	}{
		"12345678909876543X",
	}
	v = New(s)
	assert.Equal(t, v.Check(), true)
	assert.Equal(t, len(v.Errors), 0)
}
