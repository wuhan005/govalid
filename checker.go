package govalid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/language"
)

// CheckFunc is the type of checker function.
// It returns an error context if error occurred.
type CheckFunc func(ctx CheckerContext) *ErrContext

// CheckerContext is the context of checker,
// which can contains the rule of the checker and the value of the current struct field.
type CheckerContext struct {
	StructValue      reflect.Value
	FieldName        string
	FieldType        reflect.Type
	FieldValue       interface{}
	FieldLabel       string
	TemplateLanguage language.Tag

	Rule *rule
}

// Checkers is the function list of checkers.
var Checkers = map[string]CheckFunc{
	"required":     required,
	"min":          min,
	"max":          max,
	"minlen":       minlen,
	"maxlen":       maxlen,
	"alpha":        alpha,
	"alphanumeric": alphaNumeric,
	"alphadash":    alphaDash,
	"username":     userName,
	"email":        email,
	"ipv4":         ipv4,
	"mobile":       mobile,
	"tel":          tel,
	"phone":        phone,
	"idcard":       idCard,
	"equal":        equal,
	"list":         list,
}

func required(c CheckerContext) *ErrContext {
	errCtx := NewErrorContext(c)

	if c.FieldValue == nil {
		return errCtx
	}

	// Length-aware kinds (slice/array/map/string/chan) are considered empty
	// when their length is zero. Doing this before the zero-value comparison
	// also avoids "comparing uncomparable type" panics for maps and slices.
	switch c.FieldType.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		if reflect.ValueOf(c.FieldValue).Len() == 0 {
			return errCtx
		}
		return nil
	case reflect.Ptr, reflect.Interface, reflect.Func:
		if reflect.ValueOf(c.FieldValue).IsNil() {
			return errCtx
		}
		return nil
	}

	// Skip incomparable types to keep reflect.Value.Interface() == FieldValue
	// from panicking.
	if !c.FieldType.Comparable() {
		return nil
	}

	zeroValue := reflect.Zero(c.FieldType)
	if zeroValue.Interface() == c.FieldValue {
		return errCtx
	}
	return nil
}

func min(c CheckerContext) *ErrContext {
	return minOrMax(c, "min")
}

func max(c CheckerContext) *ErrContext {
	return minOrMax(c, "max")
}

func minOrMax(c CheckerContext, flag string) *ErrContext {
	ctx := NewErrorContext(c)
	if len(c.Rule.params) != 1 {
		return MakeCheckerParamError(c)
	}

	if c.FieldValue == nil {
		return MakeValueTypeError(c)
	}

	limitStr := c.Rule.params[0]

	switch reflect.TypeOf(c.FieldValue).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return MakeCheckerParamError(c)
		}
		ctx.SetFieldLimitValue(limit)

		value := reflect.ValueOf(c.FieldValue).Int()
		if flag == "min" {
			if value < limit {
				return ctx
			}
		} else {
			if value > limit {
				return ctx
			}
		}
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		limit, err := strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return MakeCheckerParamError(c)
		}
		ctx.SetFieldLimitValue(limit)

		value := reflect.ValueOf(c.FieldValue).Uint()
		if flag == "min" {
			if value < limit {
				return ctx
			}
		} else {
			if value > limit {
				return ctx
			}
		}
		return nil

	case reflect.Float32, reflect.Float64:
		limit, err := strconv.ParseFloat(limitStr, 64)
		if err != nil {
			return MakeCheckerParamError(c)
		}
		ctx.SetFieldLimitValue(limit)

		value64, ok := c.FieldValue.(float64)
		if !ok {
			value32, ok := c.FieldValue.(float32)
			if !ok {
				return MakeValueTypeError(c)
			}
			value64 = float64(value32)
		}

		if flag == "min" {
			if value64 < limit {
				return ctx
			}
		} else {
			if value64 > limit {
				return ctx
			}
		}
		return nil
	}
	return nil
}

func minlen(c CheckerContext) *ErrContext {
	return minOrMaxLen(c, "min")
}

func maxlen(c CheckerContext) *ErrContext {
	return minOrMaxLen(c, "max")
}

func minOrMaxLen(c CheckerContext, flag string) *ErrContext {
	ctx := NewErrorContext(c)
	if len(c.Rule.params) != 1 {
		return MakeCheckerParamError(c)
	}

	limit, err := strconv.ParseInt(c.Rule.params[0], 10, 64)
	if err != nil {
		return MakeCheckerParamError(c)
	}
	ctx.SetFieldLimitValue(limit)

	if c.FieldValue == nil {
		return MakeValueTypeError(c)
	}

	var length int
	switch reflect.TypeOf(c.FieldValue).Kind() {
	case reflect.String:
		s := c.FieldValue.(string)
		// Empty strings short-circuit to "no error" for parity with other
		// string-aware checkers (alphaDash/email/...).
		if s == "" {
			return nil
		}
		length = utf8.RuneCountInString(s)
	case reflect.Slice, reflect.Array, reflect.Map:
		length = reflect.ValueOf(c.FieldValue).Len()
	default:
		return MakeValueTypeError(c)
	}

	if flag == "min" {
		if int64(length) < limit {
			return ctx
		}
	} else {
		if int64(length) > limit {
			return ctx
		}
	}
	return nil
}

func alpha(c CheckerContext) *ErrContext {
	if c.FieldValue == nil || reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}

	ctx := NewErrorContext(c)
	for _, v := range value {
		if v < 'A' || (v > 'Z' && v < 'a') || v > 'z' {
			return ctx
		}
	}
	return nil
}

func alphaNumeric(c CheckerContext) *ErrContext {
	if c.FieldValue == nil || reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}

	ctx := NewErrorContext(c)
	for _, v := range value {
		if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
			return ctx
		}
	}
	return nil
}

var alphaDashPattern = regexp.MustCompile(`^\w+$`)

func alphaDash(c CheckerContext) *ErrContext {
	if c.FieldValue == nil || reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}
	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}

	if !alphaDashPattern.MatchString(value) {
		return NewErrorContext(c)
	}
	return nil
}

func userName(c CheckerContext) *ErrContext {
	if c.FieldValue == nil || reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}

	// is alpha dash.
	if ctx := alphaDash(c); ctx != nil {
		ctx.SetTemplate("alphadash")
		return ctx
	}

	// first char must be a alpha. Use the first rune (not byte) so multi-byte
	// characters are handled correctly.
	firstRune, _ := utf8.DecodeRuneInString(value)
	tmp := c
	tmp.FieldValue = string(firstRune)
	if ctx := alpha(tmp); ctx != nil {
		ctx.SetTemplate("firstCharAlpha")
		return ctx
	}

	// last char can't be dash.
	if strings.HasSuffix(value, "_") {
		ctx := NewErrorContext(c)
		ctx.SetTemplate("lastUnderline")
		return ctx
	}
	return nil
}

// emailPattern is compiled once on package init to avoid re-compiling on
// every email() invocation.
var emailPattern = regexp.MustCompile(`^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`)

func email(c CheckerContext) *ErrContext {
	if c.FieldValue == nil || reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}
	if !emailPattern.MatchString(value) {
		return NewErrorContext(c)
	}
	return nil
}

var ipv4Pattern = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)

func ipv4(c CheckerContext) *ErrContext {
	if c.FieldValue == nil || reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}
	if !ipv4Pattern.MatchString(value) {
		return NewErrorContext(c)
	}
	return nil
}

// MobilePattern is used to check the mobile phone.
// Refer to https://github.com/VincentSit/ChinaMobilePhoneNumberRegex
var MobilePattern = regexp.MustCompile(`^(?:\+?86)?1(?:3\d{3}|5[^4\D]\d{2}|8\d{3}|7(?:[0-35-9]\d{2}|4(?:0\d|1[0-2]|9\d))|9[0-35-9]\d{2}|6[2567]\d{2}|4(?:(?:10|4[01])\d{3}|[68]\d{4}|[579]\d{2}))\d{6}$`)

func mobile(c CheckerContext) *ErrContext {
	if c.FieldValue == nil || reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}
	if !MobilePattern.MatchString(value) {
		return NewErrorContext(c)
	}
	return nil
}

var telPattern = regexp.MustCompile(`^(0\d{2,3}(-)?)?\d{7,8}$`)

func tel(c CheckerContext) *ErrContext {
	if c.FieldValue == nil || reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}
	if !telPattern.MatchString(value) {
		return NewErrorContext(c)
	}
	return nil
}

func phone(c CheckerContext) *ErrContext {
	telErrCtx := tel(c)
	mobileErrCtx := mobile(c)
	if telErrCtx == nil || mobileErrCtx == nil {
		return nil
	}

	// Both forms failed; surface a single, dedicated "phone" error so the
	// message reads naturally instead of mentioning either tel or mobile.
	ctx := NewErrorContext(c)
	return ctx
}

var idCardPattern = regexp.MustCompile(`(^\d{15}$)|(^\d{17}([0-9Xx])$)`)

func idCard(c CheckerContext) *ErrContext {
	if c.FieldValue == nil || reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}
	if !idCardPattern.MatchString(value) {
		return NewErrorContext(c)
	}
	return nil
}

func equal(c CheckerContext) *ErrContext {
	if len(c.Rule.params) != 1 || c.Rule.params[0] == "" {
		return MakeCheckerParamError(c)
	}

	// equal needs the surrounding struct to compare another field. If a
	// caller invokes the checker directly with a zero StructValue (or a
	// non-struct value), there's nothing to compare to.
	if !c.StructValue.IsValid() || c.StructValue.Kind() != reflect.Struct {
		return MakeFieldNotFoundError(c)
	}

	ctx := NewErrorContext(c)
	value := fmt.Sprintf("%v", c.FieldValue)
	equalField := c.Rule.params[0]

	structType := c.StructValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		if structType.Field(i).Name == equalField {
			equalFieldValue := fmt.Sprintf("%v", c.StructValue.Field(i).Interface())
			if value != equalFieldValue {
				return ctx
			}
			return nil
		}
	}
	return MakeFieldNotFoundError(c)
}

func list(c CheckerContext) *ErrContext {
	// `valid:"list:"` is parsed as a single empty param. Treat both the
	// empty slice and a slice of only empty strings as a missing-param
	// error so callers can't accidentally allow only the empty value.
	if len(c.Rule.params) == 0 {
		return MakeCheckerParamError(c)
	}
	hasNonEmpty := false
	for _, p := range c.Rule.params {
		if p != "" {
			hasNonEmpty = true
			break
		}
	}
	if !hasNonEmpty {
		return MakeCheckerParamError(c)
	}

	ctx := NewErrorContext(c)
	value := fmt.Sprintf("%v", c.FieldValue)
	for _, v := range c.Rule.params {
		if value == v {
			return nil
		}
	}
	return ctx
}
