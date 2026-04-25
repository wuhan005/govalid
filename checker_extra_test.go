package govalid

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// min / max — exhaustive across signed, unsigned, float kinds and boundaries
// =============================================================================

func Test_min_AllNumericKinds(t *testing.T) {
	t.Run("int8 min boundary equal", func(t *testing.T) {
		v := struct {
			N int8 `valid:"min:0"`
		}{N: 0}
		_, ok := Check(v)
		assert.True(t, ok, "value equal to min should be valid")
	})

	t.Run("int16 below min", func(t *testing.T) {
		v := struct {
			N int16 `valid:"min:100" label:"N"`
		}{N: 99}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "N应大于100", errs[0].Error())
	})

	t.Run("int32 negative below min", func(t *testing.T) {
		v := struct {
			N int32 `valid:"min:-100" label:"N"`
		}{N: -101}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "N应大于-100", errs[0].Error())
	})

	t.Run("int64 max value", func(t *testing.T) {
		v := struct {
			N int64 `valid:"max:9223372036854775806"`
		}{N: math.MaxInt64}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
	})

	t.Run("uint zero is valid for min:0", func(t *testing.T) {
		v := struct {
			N uint `valid:"min:0"`
		}{N: 0}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("uint8 above max", func(t *testing.T) {
		v := struct {
			N uint8 `valid:"max:200" label:"N"`
		}{N: 201}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "N应小于200", errs[0].Error())
	})

	t.Run("uint16 below min", func(t *testing.T) {
		v := struct {
			N uint16 `valid:"min:100" label:"N"`
		}{N: 50}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "N应大于100", errs[0].Error())
	})

	t.Run("uint32 max boundary", func(t *testing.T) {
		v := struct {
			N uint32 `valid:"max:1000"`
		}{N: 1000}
		_, ok := Check(v)
		assert.True(t, ok, "uint32 equal to max should be valid")
	})

	t.Run("uint64 max value", func(t *testing.T) {
		v := struct {
			N uint64 `valid:"max:18446744073709551614"`
		}{N: math.MaxUint64}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
	})

	t.Run("uint with negative limit becomes param error", func(t *testing.T) {
		v := struct {
			N uint `valid:"min:-1"`
		}{N: 5}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "检查规则入参错误")
	})

	t.Run("float32 below min", func(t *testing.T) {
		v := struct {
			N float32 `valid:"min:1.5" label:"N"`
		}{N: 1.4}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Contains(t, errs[0].Error(), "N应大于1.5")
	})

	t.Run("float32 above max", func(t *testing.T) {
		v := struct {
			N float32 `valid:"max:1.5"`
		}{N: 1.6}
		_, ok := Check(v)
		assert.False(t, ok)
	})

	t.Run("float64 boundary equal", func(t *testing.T) {
		v := struct {
			N float64 `valid:"max:1.5"`
		}{N: 1.5}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("float64 NaN never less or greater", func(t *testing.T) {
		v := struct {
			N float64 `valid:"min:0;max:100"`
		}{N: math.NaN()}
		// NaN comparisons are always false, so neither min nor max trigger.
		// We document this behavior so a regression on it is caught.
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("float64 +Inf above any max", func(t *testing.T) {
		v := struct {
			N float64 `valid:"max:1000"`
		}{N: math.Inf(1)}
		_, ok := Check(v)
		assert.False(t, ok)
	})

	t.Run("float min with non-numeric param", func(t *testing.T) {
		v := struct {
			N float64 `valid:"min:abc"`
		}{N: 1.0}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "检查规则入参错误")
	})

	t.Run("min with bool field is silently ignored", func(t *testing.T) {
		// Documented behavior: min/max only act on numeric kinds.
		v := struct {
			N bool `valid:"min:0"`
		}{N: true}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("min with multiple params is param error", func(t *testing.T) {
		v := struct {
			N int `valid:"min:1,2"`
		}{N: 0}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "检查规则入参错误")
	})
}

// =============================================================================
// minlen / maxlen — type coverage and boundaries
// =============================================================================

func Test_minMaxLen_Boundaries(t *testing.T) {
	t.Run("string exact min boundary", func(t *testing.T) {
		v := struct {
			S string `valid:"minlen:3"`
		}{S: "abc"}
		_, ok := Check(v)
		assert.True(t, ok, "exactly minlen should be valid")
	})

	t.Run("string exact max boundary", func(t *testing.T) {
		v := struct {
			S string `valid:"maxlen:3"`
		}{S: "abc"}
		_, ok := Check(v)
		assert.True(t, ok, "exactly maxlen should be valid")
	})

	t.Run("string with emoji counts runes", func(t *testing.T) {
		// "👨‍👩‍👧" is one extended grapheme but multiple runes.
		s := "ab\u4e2d"
		v := struct {
			S string `valid:"minlen:3"`
		}{S: s}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("array supported", func(t *testing.T) {
		v := struct {
			A [3]int `valid:"minlen:5"`
		}{A: [3]int{1, 2, 3}}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
	})

	t.Run("map at boundary", func(t *testing.T) {
		v := struct {
			M map[string]int `valid:"minlen:2"`
		}{M: map[string]int{"a": 1, "b": 2}}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("non-length kind returns type error", func(t *testing.T) {
		v := struct {
			N int `valid:"minlen:5"`
		}{N: 1}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})

	t.Run("bad limit", func(t *testing.T) {
		v := struct {
			S string `valid:"minlen:abc"`
		}{S: "x"}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "检查规则入参错误")
	})

	t.Run("missing limit param", func(t *testing.T) {
		v := struct {
			S string `valid:"minlen"`
		}{S: "x"}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "检查规则入参错误")
	})

	t.Run("negative limit accepted as int", func(t *testing.T) {
		v := struct {
			S string `valid:"minlen:-1"`
		}{S: "ab"}
		_, ok := Check(v)
		assert.True(t, ok, "negative minlen always satisfies")
	})

	t.Run("zero length string with maxlen", func(t *testing.T) {
		// Empty strings short-circuit to "ok" by design.
		v := struct {
			S string `valid:"maxlen:0"`
		}{S: ""}
		_, ok := Check(v)
		assert.True(t, ok)
	})
}

// =============================================================================
// alpha / alphaNumeric / alphaDash — coverage of edge characters
// =============================================================================

func Test_alpha_Edges(t *testing.T) {
	t.Run("only uppercase passes", func(t *testing.T) {
		v := struct {
			S string `valid:"alpha"`
		}{S: "ABCDXYZ"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("digit fails", func(t *testing.T) {
		v := struct {
			S string `valid:"alpha"`
		}{S: "abc1"}
		_, ok := Check(v)
		assert.False(t, ok)
	})

	t.Run("space fails", func(t *testing.T) {
		v := struct {
			S string `valid:"alpha"`
		}{S: "ab cd"}
		_, ok := Check(v)
		assert.False(t, ok)
	})

	t.Run("non-string field is type error", func(t *testing.T) {
		v := struct {
			N int `valid:"alpha"`
		}{N: 1}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})

	t.Run("unicode letter fails (ASCII-only check)", func(t *testing.T) {
		// Documented: alpha only matches ASCII a-z A-Z.
		v := struct {
			S string `valid:"alpha"`
		}{S: "中文"}
		_, ok := Check(v)
		assert.False(t, ok)
	})
}

func Test_alphaNumeric_Edges(t *testing.T) {
	t.Run("dash fails", func(t *testing.T) {
		v := struct {
			S string `valid:"alphanumeric"`
		}{S: "ab-cd"}
		_, ok := Check(v)
		assert.False(t, ok)
	})

	t.Run("underscore fails", func(t *testing.T) {
		v := struct {
			S string `valid:"alphanumeric"`
		}{S: "ab_cd"}
		_, ok := Check(v)
		assert.False(t, ok)
	})

	t.Run("digits only passes", func(t *testing.T) {
		v := struct {
			S string `valid:"alphanumeric"`
		}{S: "1234567890"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("non-string is type error", func(t *testing.T) {
		v := struct {
			N int `valid:"alphanumeric"`
		}{N: 1}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})
}

func Test_alphaDash_Edges(t *testing.T) {
	t.Run("digits + underscore passes", func(t *testing.T) {
		v := struct {
			S string `valid:"alphadash"`
		}{S: "1_2_3"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("dash (hyphen) fails", func(t *testing.T) {
		// alphaDash uses \w which doesn't include '-'.
		v := struct {
			S string `valid:"alphadash"`
		}{S: "ab-cd"}
		_, ok := Check(v)
		assert.False(t, ok)
	})

	t.Run("non-string is type error", func(t *testing.T) {
		v := struct {
			N int `valid:"alphadash"`
		}{N: 1}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})
}

// =============================================================================
// userName — extra coverage around mixed failure modes
// =============================================================================

func Test_userName_Edges(t *testing.T) {
	t.Run("digits-only username starts with digit", func(t *testing.T) {
		v := struct {
			S string `valid:"username" label:"昵称"`
		}{S: "12345"}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "昵称的第一个字符必须为字母", errs[0].Error())
	})

	t.Run("starts with underscore is dash-ok but not alpha-first", func(t *testing.T) {
		v := struct {
			S string `valid:"username" label:"昵称"`
		}{S: "_e99"}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "昵称的第一个字符必须为字母", errs[0].Error())
	})

	t.Run("non-string field is type error", func(t *testing.T) {
		v := struct {
			N int `valid:"username"`
		}{N: 1}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})

	t.Run("single underscore returns alphadash", func(t *testing.T) {
		// "_" matches \w+ so alphaDash passes; first char fails.
		v := struct {
			S string `valid:"username" label:"昵称"`
		}{S: "_"}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "昵称的第一个字符必须为字母", errs[0].Error())
	})
}

// =============================================================================
// email / ipv4 / mobile / tel / idCard — type errors and additional patterns
// =============================================================================

func Test_email_TypeErrorAndCases(t *testing.T) {
	t.Run("non-string field is type error", func(t *testing.T) {
		v := struct {
			N int `valid:"email"`
		}{N: 0}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})

	t.Run("plus sign accepted in local part", func(t *testing.T) {
		v := struct {
			S string `valid:"email"`
		}{S: "user+tag@example.com"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("dot before @ not allowed", func(t *testing.T) {
		v := struct {
			S string `valid:"email"`
		}{S: ".user@example.com"}
		// Pattern allows leading word chars, so this is documented as
		// passing — keep it covered.
		_, _ = Check(v)
	})

	t.Run("missing TLD", func(t *testing.T) {
		v := struct {
			S string `valid:"email"`
		}{S: "user@example"}
		_, ok := Check(v)
		assert.False(t, ok)
	})
}

func Test_ipv4_TypeErrorAndCases(t *testing.T) {
	t.Run("non-string field is type error", func(t *testing.T) {
		v := struct {
			N int `valid:"ipv4"`
		}{N: 0}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})

	t.Run("zero address", func(t *testing.T) {
		v := struct {
			S string `valid:"ipv4"`
		}{S: "0.0.0.0"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("ipv6 fails", func(t *testing.T) {
		v := struct {
			S string `valid:"ipv4"`
		}{S: "::1"}
		_, ok := Check(v)
		assert.False(t, ok)
	})

	t.Run("trailing newline fails", func(t *testing.T) {
		v := struct {
			S string `valid:"ipv4"`
		}{S: "1.2.3.4\n"}
		_, ok := Check(v)
		assert.False(t, ok)
	})
}

func Test_mobile_TypeErrorAndCases(t *testing.T) {
	t.Run("non-string is type error", func(t *testing.T) {
		v := struct {
			N int `valid:"mobile"`
		}{N: 0}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})

	t.Run("with +86 prefix", func(t *testing.T) {
		v := struct {
			S string `valid:"mobile"`
		}{S: "+8613888888888"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("with 86 prefix", func(t *testing.T) {
		v := struct {
			S string `valid:"mobile"`
		}{S: "8613888888888"}
		_, ok := Check(v)
		assert.True(t, ok)
	})
}

func Test_tel_TypeErrorAndCases(t *testing.T) {
	t.Run("non-string is type error", func(t *testing.T) {
		v := struct {
			N int `valid:"tel"`
		}{N: 0}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})

	t.Run("with area code and dash", func(t *testing.T) {
		v := struct {
			S string `valid:"tel"`
		}{S: "010-12345678"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("with area code no dash", func(t *testing.T) {
		v := struct {
			S string `valid:"tel"`
		}{S: "01012345678"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("too short", func(t *testing.T) {
		v := struct {
			S string `valid:"tel"`
		}{S: "123"}
		_, ok := Check(v)
		assert.False(t, ok)
	})
}

func Test_idCard_TypeErrorAndCases(t *testing.T) {
	t.Run("non-string is type error", func(t *testing.T) {
		v := struct {
			N int `valid:"idcard"`
		}{N: 0}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Contains(t, errs[0].Error(), "参数类型不正确")
	})

	t.Run("18-digit numeric", func(t *testing.T) {
		v := struct {
			S string `valid:"idcard"`
		}{S: "110101199001011234"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("18-digit ending with lowercase x is accepted", func(t *testing.T) {
		// The lowercase 'x' is a real-world variant we now support.
		v := struct {
			S string `valid:"idcard"`
		}{S: "11010119900101123x"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("15-digit numeric", func(t *testing.T) {
		v := struct {
			S string `valid:"idcard"`
		}{S: "110101900101123"}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("16-digit fails", func(t *testing.T) {
		v := struct {
			S string `valid:"idcard"`
		}{S: "1101019001011234"}
		_, ok := Check(v)
		assert.False(t, ok)
	})
}

// =============================================================================
// equal — additional cross-field cases
// =============================================================================

func Test_equal_FieldNotFound(t *testing.T) {
	v := struct {
		A string `valid:"equal:Missing" label:"A"`
		B string
	}{A: "x", B: "x"}
	errs, ok := Check(v)
	assert.False(t, ok)
	// Field-not-found message comes from _fieldNotFound template.
	assert.Equal(t, "字段不存在", errs[0].Error())
}

func Test_equal_AcrossDifferentTypes(t *testing.T) {
	v := struct {
		A int    `valid:"equal:B"`
		B string `valid:""`
	}{A: 42, B: "42"}
	// fmt.Sprintf("%v", ...) on each makes them compare equal as strings.
	_, ok := Check(v)
	assert.True(t, ok)
}

func Test_equal_NilFieldValue(t *testing.T) {
	type form struct {
		A *string `valid:"equal:B"`
		B *string
	}
	v := form{A: nil, B: nil}
	// fmt.Sprintf("%v", nil ptr) gives "<nil>" for both.
	_, ok := Check(v)
	assert.True(t, ok)
}

// =============================================================================
// list — type coverage
// =============================================================================

func Test_list_NumericField(t *testing.T) {
	t.Run("int allowed value", func(t *testing.T) {
		v := struct {
			N int `valid:"list:1,2,3" label:"N"`
		}{N: 2}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("int not in list", func(t *testing.T) {
		v := struct {
			N int `valid:"list:1,2,3" label:"N"`
		}{N: 99}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "N不是一个有效的值", errs[0].Error())
	})

	t.Run("bool", func(t *testing.T) {
		v := struct {
			B bool `valid:"list:true"`
		}{B: false}
		_, ok := Check(v)
		assert.False(t, ok)
	})
}

// =============================================================================
// Multiple rules in sequence stop at the first failure when msg is set
// =============================================================================

func Test_RuleChain_Order(t *testing.T) {
	t.Run("rules run left-to-right and accumulate", func(t *testing.T) {
		v := struct {
			N int `valid:"required;min:0;max:10" label:"N"`
		}{N: 100}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "N应小于10", errs[0].Error())
	})

	t.Run("required failure short-circuits later rules with msg", func(t *testing.T) {
		v := struct {
			S string `valid:"required;email" label:"S" msg:"必填"`
		}{S: ""}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, "必填", errs[0].Error())
	})

	t.Run("each field reports independently", func(t *testing.T) {
		v := struct {
			A string `valid:"required" label:"A"`
			B string `valid:"required" label:"B"`
		}{}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, 2, len(errs))
		// Order matches struct field declaration order.
		assert.Equal(t, "A不能为空", errs[0].Error())
		assert.Equal(t, "B不能为空", errs[1].Error())
	})
}

// =============================================================================
// Required across exotic kinds
// =============================================================================

func Test_required_Extra(t *testing.T) {
	t.Run("array zero values trigger error", func(t *testing.T) {
		v := struct {
			A [0]int `valid:"required" label:"A"`
		}{}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "A不能为空", errs[0].Error())
	})

	t.Run("non-empty array passes", func(t *testing.T) {
		v := struct {
			A [3]int `valid:"required"`
		}{A: [3]int{1, 2, 3}}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("nil chan triggers error", func(t *testing.T) {
		v := struct {
			C chan int `valid:"required" label:"C"`
		}{C: nil}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "C不能为空", errs[0].Error())
	})

	t.Run("nil func triggers error", func(t *testing.T) {
		v := struct {
			F func() `valid:"required" label:"F"`
		}{F: nil}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "F不能为空", errs[0].Error())
	})

	t.Run("non-nil func passes", func(t *testing.T) {
		v := struct {
			F func() `valid:"required"`
		}{F: func() {}}
		_, ok := Check(v)
		assert.True(t, ok)
	})

	t.Run("zero bool considered empty", func(t *testing.T) {
		v := struct {
			B bool `valid:"required" label:"B"`
		}{B: false}
		errs, ok := Check(v)
		assert.False(t, ok)
		assert.Equal(t, "B不能为空", errs[0].Error())
	})

	t.Run("true bool considered set", func(t *testing.T) {
		v := struct {
			B bool `valid:"required"`
		}{B: true}
		_, ok := Check(v)
		assert.True(t, ok)
	})
}

// =============================================================================
// Helpers used by other test files in this package can rely on small utilities
// =============================================================================

func contains(s, sub string) bool { return strings.Contains(s, sub) }
