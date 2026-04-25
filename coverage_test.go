package govalid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

// =============================================================================
// SetMessageTemplates: registering a brand-new locale that wasn't
// pre-populated by the package.
// =============================================================================

func Test_SetMessageTemplates_NewLocale(t *testing.T) {
	// Pick a locale that the package doesn't ship templates for.
	loc := language.Japanese
	defer delete(errorTemplateSet, loc)

	SetMessageTemplates(map[string]string{
		"required": " は必須です",
	}, loc)

	v := struct {
		Name string `valid:"required" label:"Name"`
	}{}
	errs, ok := Check(v, loc)
	assert.False(t, ok)
	assert.Equal(t, "Name は必須です", errs[0].Error())
}

// =============================================================================
// minOrMaxLen: float field is rejected with a type error (default branch).
// =============================================================================

func Test_minlen_FloatTypeError(t *testing.T) {
	v := struct {
		N float64 `valid:"minlen:5"`
	}{N: 1.0}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Contains(t, errs[0].Error(), "参数类型不正确")
}

// =============================================================================
// makeMessage: template that contains placeholders is rendered with the
// limit value substituted.
// =============================================================================

func Test_makeMessage_Placeholders(t *testing.T) {
	originalCN := errorTemplateChinese["min"]
	defer func() { errorTemplateChinese["min"] = originalCN }()

	// Override the template to use placeholders so we hit the
	// strings.NewReplacer branch in makeMessage.
	SetMessageTemplates(map[string]string{
		"min": "{{字段 {field} 必须 ≥ {limit}}}",
	})

	v := struct {
		Score int `valid:"min:10" label:"分数"`
	}{Score: 5}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	// "{field}" is replaced with the field's *programmatic* name
	// (Score), and {limit} with the parsed limit.
	assert.Equal(t, "字段 Score 必须 ≥ 10", errs[0].Error())
}

// =============================================================================
// equal: invoking the checker directly with a non-struct StructValue
// surfaces a field-not-found error rather than panicking.
// =============================================================================

func Test_equal_DirectInvocation_BadStruct(t *testing.T) {
	// Direct synthetic call mimicking what a third-party caller could do.
	out := equal(CheckerContext{
		FieldName:        "X",
		FieldLabel:       "X",
		FieldValue:       "v",
		FieldType:        reflectStringType,
		TemplateLanguage: language.Chinese,
		Rule:             &rule{checker: "equal", params: []string{"Other"}},
	})
	assert.NotNil(t, out)
	assert.Equal(t, "字段不存在", out.Error())
}

// =============================================================================
// list directly invoked with a nil FieldValue: stringifies to "<nil>"
// which won't match anything.
// =============================================================================

func Test_list_NilFieldValue_DoesNotMatch(t *testing.T) {
	out := list(CheckerContext{
		FieldName:        "X",
		FieldLabel:       "X",
		FieldValue:       nil,
		FieldType:        reflectStringType,
		TemplateLanguage: language.Chinese,
		Rule:             &rule{checker: "list", params: []string{"a", "b"}},
	})
	assert.NotNil(t, out)
}

// =============================================================================
// required: directly invoked with a nil FieldValue should still trigger
// the required error.
// =============================================================================

func Test_required_DirectInvocation_NilValue(t *testing.T) {
	out := required(CheckerContext{
		FieldName:        "X",
		FieldLabel:       "X",
		FieldValue:       nil,
		FieldType:        reflectStringType,
		TemplateLanguage: language.Chinese,
		Rule:             &rule{checker: "required"},
	})
	assert.NotNil(t, out)
	assert.Equal(t, "X不能为空", out.Error())
}
