package govalid

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

// =============================================================================
// Realistic end-to-end: a registration form covering most checkers, msg
// override, equal cross-field, custom Validate(), and a struct slice.
// =============================================================================

type address struct {
	Street string `valid:"required;maxlen:200" label:"街道"`
	Zip    string `valid:"required;alphanumeric;minlen:5;maxlen:10" label:"邮编"`
}

type registerForm struct {
	Username        string    `valid:"required;username;minlen:3;maxlen:20" label:"用户名"`
	Email           string    `valid:"required;email" label:"邮箱"`
	Mobile          string    `valid:"required;mobile" label:"手机号"`
	Age             uint      `valid:"required;min:18;max:120" label:"年龄"`
	Score           float64   `valid:"min:0;max:100" label:"分数"`
	Password        string    `valid:"required;minlen:8;maxlen:64" label:"密码" msg:"密码长度需在 8-64 之间"`
	ConfirmPassword string    `valid:"required;equal:Password" label:"确认密码"`
	Role            string    `valid:"list:admin,editor,viewer" label:"角色"`
	IDCard          string    `valid:"idcard" label:"身份证"`
	HomeIPv4        string    `valid:"ipv4" label:"家庭 IP"`
	Tags            []string  `valid:"required;minlen:1" label:"标签"`
	Addresses       []address `valid:"required" label:"地址"`
}

func (f *registerForm) Validate() error {
	if f.Username == "admin" {
		return errors.New("不允许使用 admin 作为用户名")
	}
	return nil
}

func Test_Integration_AllValid(t *testing.T) {
	f := &registerForm{
		Username:        "iwh",
		Email:           "i@example.com",
		Mobile:          "13888888888",
		Age:             24,
		Score:           99.5,
		Password:        "supersecret",
		ConfirmPassword: "supersecret",
		Role:            "admin",
		IDCard:          "123456789012345",
		HomeIPv4:        "127.0.0.1",
		Tags:            []string{"tag1"},
		Addresses: []address{
			{Street: "无名街 1 号", Zip: "100000"},
		},
	}
	// admin name is rejected by Validate but everything else passes; run
	// with a non-admin name.
	f.Username = "user"
	errs, ok := Check(f)
	assert.True(t, ok, "errs=%v", errs)
	assert.Empty(t, errs)
}

func Test_Integration_AllInvalid(t *testing.T) {
	f := &registerForm{
		Username:        "1bad",
		Email:           "no-at-sign",
		Mobile:          "12345",
		Age:             5,
		Score:           1000,
		Password:        "short",
		ConfirmPassword: "different",
		Role:            "wizard",
		IDCard:          "abc",
		HomeIPv4:        "999.999.999.999",
		Tags:            nil,
		Addresses:       nil,
	}
	errs, ok := Check(f)
	assert.False(t, ok)
	// We expect at least one error per field (12 + 1 nested validate failure
	// blocked because name failed first; "admin"-block only triggers when
	// username == "admin"). We just check we got many errors, and the first
	// few are in field-declaration order.
	assert.GreaterOrEqual(t, len(errs), 10)
	assert.Equal(t, "用户名的第一个字符必须为字母", errs[0].Error())
}

func Test_Integration_AdminBlocked(t *testing.T) {
	f := &registerForm{
		Username:        "admin",
		Email:           "i@example.com",
		Mobile:          "13888888888",
		Age:             24,
		Password:        "supersecret",
		ConfirmPassword: "supersecret",
		Role:            "admin",
		Tags:            []string{"tag1"},
		Addresses: []address{
			{Street: "S", Zip: "100000"},
		},
	}
	errs, ok := Check(f)
	assert.False(t, ok)

	hasAdminBlock := false
	for _, e := range errs {
		if e.Error() == "不允许使用 admin 作为用户名" {
			hasAdminBlock = true
		}
	}
	assert.True(t, hasAdminBlock, "Validate() error should appear")
}

// =============================================================================
// Default vs explicit language fallback through Check
// =============================================================================

func Test_Integration_LanguageFallback(t *testing.T) {
	v := struct {
		Name string `valid:"required" label:"用户名" label-en:"username"`
	}{}

	t.Run("default uses chinese", func(t *testing.T) {
		errs, _ := Check(v)
		assert.Equal(t, "用户名不能为空", errs[0].Error())
	})

	t.Run("english picks english", func(t *testing.T) {
		errs, _ := Check(v, language.English)
		assert.Equal(t, "username can not be empty", errs[0].Error())
	})

	t.Run("unknown locale falls back to chinese template", func(t *testing.T) {
		errs, _ := Check(v, language.Korean)
		// The template falls back to chinese, but the label still picks
		// the unqualified "label" tag because there's no label-ko.
		assert.Equal(t, "用户名不能为空", errs[0].Error())
	})
}

// =============================================================================
// Checkers map can be removed/replaced safely
// =============================================================================

func Test_Integration_CheckerOverride(t *testing.T) {
	// Save and restore "required" to leave a clean slate.
	original := Checkers["required"]
	defer func() { Checkers["required"] = original }()

	Checkers["required"] = func(c CheckerContext) *ErrContext {
		// Make required permissive.
		return nil
	}

	v := struct {
		S string `valid:"required" label:"S"`
	}{S: ""}
	_, ok := Check(v)
	assert.True(t, ok, "overridden required should pass")
}
