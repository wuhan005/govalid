package govalid

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

// Test_NilInput verifies that Check does not panic on nil or typed-nil
// input.
func Test_NilInput(t *testing.T) {
	t.Run("untyped nil", func(t *testing.T) {
		errs, ok := Check(nil)
		assert.True(t, ok)
		assert.Nil(t, errs)
	})

	t.Run("typed nil pointer", func(t *testing.T) {
		type form struct {
			Name string `valid:"required"`
		}
		var f *form
		errs, ok := Check(f)
		assert.True(t, ok)
		assert.Nil(t, errs)
	})

	t.Run("non-struct kind", func(t *testing.T) {
		errs, ok := Check(42)
		assert.True(t, ok)
		assert.Nil(t, errs)

		errs, ok = Check("hello")
		assert.True(t, ok)
		assert.Nil(t, errs)
	})
}

// Test_RequiredMap verifies that the required checker handles maps and
// other length-aware kinds without panicking on incomparable types.
func Test_RequiredMap(t *testing.T) {
	t.Run("nil map", func(t *testing.T) {
		v := struct {
			Data map[string]string `valid:"required" label:"配置"`
		}{}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "配置不能为空", errs[0].Error())
	})

	t.Run("empty map", func(t *testing.T) {
		v := struct {
			Data map[string]string `valid:"required" label:"配置"`
		}{Data: map[string]string{}}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
	})

	t.Run("populated map", func(t *testing.T) {
		v := struct {
			Data map[string]string `valid:"required"`
		}{Data: map[string]string{"k": "v"}}
		errs, ok := Check(v)
		assert.True(t, ok)
		assert.Nil(t, errs)
	})
}

// Test_RequiredPtr verifies that required correctly detects nil pointers.
func Test_RequiredPtr(t *testing.T) {
	type inner struct{ X int }

	t.Run("nil ptr", func(t *testing.T) {
		v := struct {
			Data *inner `valid:"required" label:"内层"`
		}{}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "内层不能为空", errs[0].Error())
	})

	t.Run("non-nil ptr", func(t *testing.T) {
		v := struct {
			Data *inner `valid:"required"`
		}{Data: &inner{}}
		errs, ok := Check(v)
		assert.True(t, ok)
		assert.Nil(t, errs)
	})
}

// Test_EqualEmptyParams verifies that equal returns a parameter error
// instead of panicking on empty params.
func Test_EqualEmptyParams(t *testing.T) {
	v := struct {
		Name string `valid:"equal" label:"姓名"`
	}{Name: "x"}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "姓名检查规则入参错误", errs[0].Error())
}

// Test_PhoneError verifies the error message for phone failures uses the
// dedicated "phone" template instead of leaking tel/mobile messages.
func Test_PhoneError(t *testing.T) {
	v := struct {
		P string `valid:"phone" label:"号码"`
	}{P: "abc"}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "号码不是合法的号码", errs[0].Error())
}

// Test_UsernameAlphaDashMessage verifies that a non-alphadash username
// produces the alphadash error message, not firstCharAlpha.
func Test_UsernameAlphaDashMessage(t *testing.T) {
	v := struct {
		Name string `valid:"username" label:"昵称"`
	}{Name: "abc$%^"}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "昵称只含有数字或字母以及下划线", errs[0].Error())
}

// Test_UsernameUnicodeFirstChar verifies that the first-char check uses
// runes (so a leading multi-byte character is still rejected cleanly).
func Test_UsernameUnicodeFirstChar(t *testing.T) {
	v := struct {
		Name string `valid:"username" label:"昵称"`
	}{Name: "用户e99"}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	// alphaDash regex \w doesn't match Chinese characters, so we expect
	// the alphadash message.
	assert.Equal(t, "昵称只含有数字或字母以及下划线", errs[0].Error())
}

// Test_ListEmptyParams verifies list returns a parameter error when no
// allowed values are provided, instead of silently rejecting non-empty
// values.
func Test_ListEmptyParams(t *testing.T) {
	t.Run("empty params", func(t *testing.T) {
		v := struct {
			Role string `valid:"list:" label:"角色"`
		}{Role: "anything"}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "角色检查规则入参错误", errs[0].Error())
	})
}

// Test_MinLenSliceMap verifies that minlen/maxlen handle slices and maps
// instead of stringifying them with %s.
func Test_MinLenSliceMap(t *testing.T) {
	t.Run("slice too short", func(t *testing.T) {
		v := struct {
			Items []int `valid:"minlen:3" label:"项目"`
		}{Items: []int{1, 2}}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "项目长度应大于3", errs[0].Error())
	})

	t.Run("slice ok", func(t *testing.T) {
		v := struct {
			Items []int `valid:"minlen:3"`
		}{Items: []int{1, 2, 3, 4}}
		errs, ok := Check(v)
		assert.True(t, ok)
		assert.Nil(t, errs)
	})

	t.Run("map too long", func(t *testing.T) {
		v := struct {
			Items map[string]int `valid:"maxlen:1" label:"映射"`
		}{Items: map[string]int{"a": 1, "b": 2}}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "映射长度应小于1", errs[0].Error())
	})
}

// Test_MinMaxNilField verifies min/max do not panic when given a nil
// interface field value.
func Test_MinMaxNilField(t *testing.T) {
	v := struct {
		Score interface{} `valid:"min:0" label:"得分"`
	}{Score: nil}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "得分参数类型不正确", errs[0].Error())
}

// Test_AlphaNilField verifies alpha/alphaNumeric/alphaDash do not panic
// on nil interface field values.
func Test_AlphaNilField(t *testing.T) {
	t.Run("alpha nil", func(t *testing.T) {
		v := struct {
			Data interface{} `valid:"alpha" label:"数据"`
		}{Data: nil}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
	})

	t.Run("alpha empty string passes", func(t *testing.T) {
		v := struct {
			Name string `valid:"alpha"`
		}{Name: ""}
		errs, ok := Check(v)
		assert.True(t, ok)
		assert.Nil(t, errs)
	})
}

// Test_SetMessageTemplates verifies that the API documented in README
// actually exists and works.
func Test_SetMessageTemplates(t *testing.T) {
	original := errorTemplateChinese["required"]
	originalEng := errorTemplateEnglish["required"]
	defer func() {
		errorTemplateChinese["required"] = original
		errorTemplateEnglish["required"] = originalEng
	}()

	t.Run("default language", func(t *testing.T) {
		SetMessageTemplates(map[string]string{
			"required": " is required",
		})

		v := struct {
			Name string `valid:"required" label:"用户名"`
		}{}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "用户名 is required", errs[0].Error())
	})

	t.Run("explicit language", func(t *testing.T) {
		SetMessageTemplates(map[string]string{
			"required": " must be present",
		}, language.English)

		v := struct {
			Name string `valid:"required" label:"User"`
		}{}
		errs, ok := Check(v, language.English)
		assert.False(t, ok)
		assert.Equal(t, "User must be present", errs[0].Error())
	})
}

// Test_ValidateMethodValueReceiver ensures that a Validate method with a
// value receiver is found on both the value and the pointer.
func Test_ValidateMethodValueReceiver(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		f := valueReceiverForm{Name: "fail"}
		errs, ok := Check(f)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "value receiver fail", errs[0].Error())
	})

	t.Run("pointer", func(t *testing.T) {
		f := &valueReceiverForm{Name: "fail"}
		errs, ok := Check(f)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "value receiver fail", errs[0].Error())
	})
}

type valueReceiverForm struct {
	Name string `valid:"required"`
}

func (v valueReceiverForm) Validate() error {
	if v.Name == "fail" {
		return errors.New("value receiver fail")
	}
	return nil
}

// Test_ValidateMethodPointerReceiver ensures that a pointer-receiver
// Validate is invoked when the caller passes a pointer.
func Test_ValidateMethodPointerReceiver(t *testing.T) {
	f := &pointerReceiverForm{Name: "fail"}
	errs, ok := Check(f)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "pointer receiver fail", errs[0].Error())
}

type pointerReceiverForm struct {
	Name string `valid:"required"`
}

func (p *pointerReceiverForm) Validate() error {
	if p.Name == "fail" {
		return errors.New("pointer receiver fail")
	}
	return nil
}

// Test_UnexportedField verifies that structs with unexported fields don't
// panic during reflection.
func Test_UnexportedField(t *testing.T) {
	type Inner struct {
		Name string `valid:"required"`
	}
	v := struct {
		inner Inner
		Name  string `valid:"required" label:"姓名"`
	}{
		inner: Inner{Name: ""},
		Name:  "",
	}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "姓名不能为空", errs[0].Error())
}
