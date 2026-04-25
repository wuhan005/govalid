package govalid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

// =============================================================================
// parseRules — additional edge cases
// =============================================================================

func Test_parseRules_Extra(t *testing.T) {
	for _, tc := range []struct {
		name string
		rule string
		want []*rule
	}{
		{
			name: "trailing semicolon",
			rule: "required;",
			want: []*rule{{checker: "required"}},
		},
		{
			name: "leading semicolon",
			rule: ";required",
			want: []*rule{{checker: "required"}},
		},
		{
			name: "double semicolon",
			rule: "required;;min:0",
			want: []*rule{
				{checker: "required"},
				{checker: "min", params: []string{"0"}},
			},
		},
		{
			name: "spaces preserved in checker name",
			// We don't trim, so callers shouldn't add spaces. Ensure that's
			// stable behavior, not silently stripping.
			rule: "required ;min:0",
			want: []*rule{
				{checker: "required "},
				{checker: "min", params: []string{"0"}},
			},
		},
		{
			name: "list with multiple values",
			rule: "list:a,b,c,d",
			want: []*rule{
				{checker: "list", params: []string{"a", "b", "c", "d"}},
			},
		},
		{
			name: "param contains colon",
			// SplitN(_, _, 2) preserves the second colon in the param.
			rule: "url:http://example.com",
			want: []*rule{
				{checker: "url", params: []string{"http://example.com"}},
			},
		},
		{
			name: "value-only with empty key skipped",
			rule: ":a,b",
			want: []*rule{},
		},
		{
			name: "all empty values",
			rule: ";;;",
			want: []*rule{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := parseRules(tc.rule)
			assert.Equal(t, tc.want, got)
		})
	}
}

// =============================================================================
// Custom RulesField / LabelField / MessageField — package-level overrides
// =============================================================================

func Test_CustomTagNames(t *testing.T) {
	// Save and restore so other tests don't see our changes.
	origR, origL, origM := RulesField, LabelField, MessageField
	defer func() {
		RulesField = origR
		LabelField = origL
		MessageField = origM
	}()

	RulesField = "v"
	LabelField = "lbl"
	MessageField = "m"

	v := struct {
		Name string `v:"required" lbl:"姓名" m:"必填~"`
	}{Name: ""}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "必填~", errs[0].Error())
}

// =============================================================================
// Language-tagged labels (label-zh, label-en)
// =============================================================================

func Test_LanguageTaggedLabel(t *testing.T) {
	type form struct {
		Name string `valid:"required" label:"DefaultLabel" label-en:"Name" label-zh:"姓名"`
	}
	v := form{Name: ""}

	t.Run("english picks label-en", func(t *testing.T) {
		errs, ok := Check(v, language.English)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "Name")
	})

	t.Run("chinese picks label-zh", func(t *testing.T) {
		errs, ok := Check(v, language.Chinese)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "姓名")
	})

	t.Run("missing locale falls back to label", func(t *testing.T) {
		// Korean isn't tagged on the field — falls back to "label".
		errs, ok := Check(v, language.Korean)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "DefaultLabel")
	})
}

// =============================================================================
// Field name fallback when no label is provided
// =============================================================================

func Test_NoLabel_UsesFieldName(t *testing.T) {
	v := struct {
		Email string `valid:"required"`
	}{}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.Equal(t, "Email不能为空", errs[0].Error())
}

// =============================================================================
// Embedded / anonymous structs
// =============================================================================

type embeddedBase struct {
	BaseField string `valid:"required" label:"基础字段"`
}

type derivedForm struct {
	embeddedBase
	Extra string `valid:"required" label:"额外字段"`
}

func Test_EmbeddedStruct(t *testing.T) {
	v := derivedForm{}
	errs, ok := Check(v)
	assert.False(t, ok)
	// Embedded fields are reachable via reflection like normal struct fields.
	assert.GreaterOrEqual(t, len(errs), 1)
	hasBase := false
	hasExtra := false
	for _, e := range errs {
		if e.Error() == "基础字段不能为空" {
			hasBase = true
		}
		if e.Error() == "额外字段不能为空" {
			hasExtra = true
		}
	}
	assert.True(t, hasBase, "embedded base field should be validated")
	assert.True(t, hasExtra, "outer field should be validated")
}

// =============================================================================
// Deeply nested structs
// =============================================================================

func Test_DeepNesting(t *testing.T) {
	type level3 struct {
		L3 string `valid:"required" label:"L3"`
	}
	type level2 struct {
		L2  string `valid:"required" label:"L2"`
		Sub level3
	}
	type level1 struct {
		L1  string `valid:"required" label:"L1"`
		Sub level2
	}

	v := level1{}
	errs, ok := Check(v)
	assert.False(t, ok)
	// Each level fires its own required error.
	msgs := make(map[string]bool)
	for _, e := range errs {
		msgs[e.Error()] = true
	}
	assert.True(t, msgs["L1不能为空"], "L1 should fail")
	assert.True(t, msgs["L2不能为空"], "L2 should fail")
	assert.True(t, msgs["L3不能为空"], "L3 should fail")
}

// =============================================================================
// Slice of struct with internal validations
// =============================================================================

func Test_SliceOfStructs_PerElementValidation(t *testing.T) {
	type item struct {
		Name string `valid:"required" label:"项目名"`
		Qty  int    `valid:"min:1" label:"数量"`
	}
	type cart struct {
		Items []item `valid:"required" label:"购物车"`
	}

	c := cart{
		Items: []item{
			{Name: "ok", Qty: 5},
			{Name: "", Qty: 0},
			{Name: "x", Qty: 10},
		},
	}
	errs, ok := Check(c)
	assert.False(t, ok)
	// Element 1 (zero-indexed) contributes two errors: name and qty.
	assert.Equal(t, 2, len(errs))
}

// =============================================================================
// Slice of pointer to struct — currently not natively supported, so we
// just document/verify behavior so a regression is caught.
// =============================================================================

func Test_SliceOfPointerToStruct_Behavior(t *testing.T) {
	type item struct {
		Name string `valid:"required" label:"项目名"`
	}
	type cart struct {
		Items []*item `valid:"required" label:"购物车"`
	}

	t.Run("non-empty slice of *item passes required (slice has length)", func(t *testing.T) {
		// parseStruct's slice-of-struct branch checks Elem().Kind() == Struct
		// for the element type, which is Ptr here, so per-element validation
		// doesn't recurse — we accept this current limitation.
		c := cart{Items: []*item{{Name: "ok"}}}
		_, ok := Check(c)
		assert.True(t, ok)
	})
}

// =============================================================================
// Anonymous unexported struct fields with their own valid tag must not
// panic on field.Interface() — the recursion still validates exported
// promoted fields.
// =============================================================================

type unexportedTaggedBase struct {
	Name string `valid:"required" label:"name"`
}

type derivedWithTaggedAnon struct {
	unexportedTaggedBase `valid:"required"`
	Other                string `valid:"required" label:"other"`
}

func Test_AnonymousUnexported_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic: %v", r)
		}
	}()
	v := derivedWithTaggedAnon{}
	errs, ok := Check(v)
	assert.False(t, ok)
	assert.NotEmpty(t, errs)
}
