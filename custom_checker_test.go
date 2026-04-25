package govalid

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Custom checker — full README example end-to-end
// =============================================================================

func Test_CustomChecker_Basic(t *testing.T) {
	// Make sure we don't pollute the global registry for other tests.
	const checkerName = "noE99"
	defer delete(Checkers, checkerName)
	defer delete(errorTemplateChinese, checkerName)
	defer delete(errorTemplateEnglish, checkerName)

	SetMessageTemplates(map[string]string{
		checkerName: "不能包含 e99",
	})

	Checkers[checkerName] = func(c CheckerContext) *ErrContext {
		v, ok := c.FieldValue.(string)
		if !ok {
			return MakeValueTypeError(c)
		}
		if strings.Contains(v, "e99") {
			return NewErrorContext(c)
		}
		return nil
	}

	t.Run("violation produces custom message", func(t *testing.T) {
		v := struct {
			Content string `valid:"noE99" label:"内容"`
		}{Content: "hello e99 world"}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "内容不能包含 e99", errs[0].Error())
	})

	t.Run("compliant value passes", func(t *testing.T) {
		v := struct {
			Content string `valid:"noE99"`
		}{Content: "hello world"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("non-string yields type error", func(t *testing.T) {
		v := struct {
			N int `valid:"noE99"`
		}{N: 0}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})
}

// =============================================================================
// Custom checker that uses Rule.params
// =============================================================================

func Test_CustomChecker_WithParams(t *testing.T) {
	const checkerName = "startsWith"
	defer delete(Checkers, checkerName)
	defer delete(errorTemplateChinese, checkerName)

	SetMessageTemplates(map[string]string{
		checkerName: "必须以指定前缀开头",
	})

	Checkers[checkerName] = func(c CheckerContext) *ErrContext {
		if len(c.Rule.params) == 0 {
			return MakeCheckerParamError(c)
		}
		v, ok := c.FieldValue.(string)
		if !ok {
			return MakeValueTypeError(c)
		}
		for _, p := range c.Rule.params {
			if strings.HasPrefix(v, p) {
				return nil
			}
		}
		return NewErrorContext(c)
	}

	t.Run("matches one of multiple prefixes", func(t *testing.T) {
		v := struct {
			S string `valid:"startsWith:foo,bar"`
		}{S: "barbaz"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("matches no prefix", func(t *testing.T) {
		v := struct {
			S string `valid:"startsWith:foo,bar" label:"S"`
		}{S: "qux"}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "S必须以指定前缀开头", errs[0].Error())
	})

	t.Run("missing param triggers param error", func(t *testing.T) {
		v := struct {
			S string `valid:"startsWith"`
		}{S: "x"}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "检查规则入参错误")
	})
}

// =============================================================================
// Validate() method — additional behavior coverage
// =============================================================================

type validateNilForm struct {
	Name string `valid:"required"`
}

func (v *validateNilForm) Validate() error { return nil }

func Test_Validate_ReturnsNil_NoError(t *testing.T) {
	f := &validateNilForm{Name: "ok"}
	errs, ok := Check(f)
	assert.True(t, ok)
	assert.Empty(t, errs)
}

type validateMultipleForm struct {
	A string `valid:"required" label:"A"`
}

func (v validateMultipleForm) Validate() error {
	return errors.New("validate said no")
}

func Test_Validate_RunsAfterTagRules(t *testing.T) {
	// When a required tag fails AND Validate fails, both errors surface.
	v := validateMultipleForm{A: ""}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 2, len(errs))
	assert.Equal(t, "A不能为空", errs[0].Error())
	assert.Equal(t, "validate said no", errs[1].Error())
}

// Validate signatures with the wrong shape are silently ignored.
type validateWrongSignature struct {
	A string `valid:"required" label:"A"`
}

func (v validateWrongSignature) Validate(extra string) error {
	return errors.New("never called")
}

func Test_Validate_WrongSignature_Ignored(t *testing.T) {
	v := validateWrongSignature{A: "x"}
	_, ok := Check(v)
	assert.True(t, ok, "Validate(string) shouldn't be invoked")
}

type validateWrongReturn struct {
	A string `valid:"required" label:"A"`
}

func (v validateWrongReturn) Validate() string {
	return "not an error"
}

func Test_Validate_WrongReturn_Ignored(t *testing.T) {
	v := validateWrongReturn{A: "x"}
	_, ok := Check(v)
	assert.True(t, ok, "Validate() string shouldn't be invoked")
}

// Custom error types satisfying the error interface should still work.
type customError struct{ msg string }

func (c *customError) Error() string { return c.msg }

type validateCustomErr struct {
	A string `valid:"required" label:"A"`
}

func (v validateCustomErr) Validate() error {
	return &customError{msg: "custom!"}
}

func Test_Validate_CustomErrorType(t *testing.T) {
	v := validateCustomErr{A: "x"}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "custom!", errs[0].Error())
}
