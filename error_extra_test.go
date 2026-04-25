package govalid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

// =============================================================================
// MakeCheckerNotFoundError — exercised when an unknown checker is referenced
// =============================================================================

func Test_MakeCheckerNotFoundError(t *testing.T) {
	t.Run("unknown checker triggers not-found error", func(t *testing.T) {
		v := struct {
			N int `valid:"thisRuleDoesNotExist" label:"N"`
		}{N: 1}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "N检查规则未找到", errs[0].Error())
	})

	t.Run("unknown checker with params", func(t *testing.T) {
		v := struct {
			N int `valid:"unknown:foo,bar"`
		}{N: 1}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "检查规则未找到")
	})

	t.Run("english locale", func(t *testing.T) {
		v := struct {
			N int `valid:"unknown" label:"N"`
		}{N: 1}
		errs, ok := Check(v, language.English)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "check rule not found")
	})
}

// =============================================================================
// MakeFieldNotFoundError — exercised by equal: pointing at a missing field
// =============================================================================

func Test_MakeFieldNotFoundError(t *testing.T) {
	v := struct {
		A string `valid:"equal:Nope"`
	}{A: "x"}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, "字段不存在", errs[0].Error())
}

// =============================================================================
// MakeUserDefinedError — direct use
// =============================================================================

func Test_MakeUserDefinedError(t *testing.T) {
	e := MakeUserDefinedError("oops")
	assert.Equal(t, "oops", e.Error())
	// Sanity check: it satisfies the error interface.
	var _ error = e
}

// =============================================================================
// ErrContext fields exposed to callers
// =============================================================================

func Test_ErrContext_PublicFields(t *testing.T) {
	v := struct {
		Name string `valid:"required" label:"姓名"`
	}{Name: ""}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))

	e := errs[0]
	assert.Equal(t, "Name", e.FieldName)
	assert.Equal(t, "姓名", e.FieldLabel)
	assert.Equal(t, "", e.FieldValue)
	assert.Equal(t, language.Chinese, e.TemplateLanguage)
}

// =============================================================================
// SetTemplate / SetFieldLimitValue mutate the rendered message
// =============================================================================

func Test_ErrContext_Mutators(t *testing.T) {
	v := struct {
		N int `valid:"min:0" label:"N"`
	}{N: -1}
	errs, _ := Check(v)
	e := errs[0]

	// Reasonable starting state.
	assert.Equal(t, "N应大于0", e.Error())

	// Override the limit value programmatically.
	e.SetFieldLimitValue(99)
	assert.Equal(t, "N应大于99", e.Error())

	// Override the template entirely.
	e.SetTemplate("required")
	assert.Equal(t, "N不能为空", e.Error())
}

// =============================================================================
// getErrorTemplate — fallbacks
// =============================================================================

func Test_getErrorTemplate_Fallbacks(t *testing.T) {
	t.Run("unknown language falls back to default chinese", func(t *testing.T) {
		// Korean isn't registered; fall back to defaultTemplateLanguage.
		got := getErrorTemplate("required", language.Korean)
		assert.Equal(t, "不能为空", got)
	})

	t.Run("unknown key falls back to _unknownErrorTemplate", func(t *testing.T) {
		got := getErrorTemplate("definitelyNotAKey", language.Chinese)
		assert.Contains(t, got, "未知错误")
	})

	t.Run("unknown key in english", func(t *testing.T) {
		got := getErrorTemplate("definitelyNotAKey", language.English)
		assert.Contains(t, got, "unknown error")
	})
}
