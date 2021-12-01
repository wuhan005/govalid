package govalid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// CheckFunc is the type of checker function.
// It returns an error context if error occurred.
type CheckFunc func(ctx CheckerContext) *ErrContext

// CheckerContext is the context of checker,
// which can contains the rule of the checker and the value of the current struct field.
type CheckerContext struct {
	FieldName  string
	FieldType  reflect.Type
	FieldValue interface{}
	FieldLabel string

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
}

func required(c CheckerContext) *ErrContext {
	errCtx := NewErrorContext(c)

	if c.FieldValue == nil {
		return errCtx
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

	limitStr := c.Rule.params[0]

	switch reflect.TypeOf(c.FieldValue).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
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

	if c.FieldValue == "" {
		return nil
	}

	value := fmt.Sprintf("%s", c.FieldValue)
	if flag == "min" {
		if int64(len(value)) < limit {
			return ctx
		}
	} else {
		if int64(len(value)) > limit {
			return ctx
		}
	}
	return nil
}

func alpha(c CheckerContext) *ErrContext {
	ctx := NewErrorContext(c)
	if reflect.ValueOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	for _, v := range c.FieldValue.(string) {
		if v < 'A' || (v > 'Z' && v < 'a') || v > 'z' {
			return ctx
		}
	}
	return nil
}

func alphaNumeric(c CheckerContext) *ErrContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	for _, v := range c.FieldValue.(string) {
		if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
			return ctx
		}
	}
	return nil
}

func alphaDash(c CheckerContext) *ErrContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}
	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}

	if !regexp.MustCompile(`^\w+$`).MatchString(value) {
		return ctx
	}
	return nil
}

func userName(c CheckerContext) *ErrContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := c.FieldValue.(string)
	if value == "" {
		return nil
	}

	// is alpha dash
	if ctx := alphaDash(c); ctx != nil {
		ctx.SetTemplate("firstCharAlpha")
		return ctx
	}

	// first char must be a alpha
	tmp := c
	tmp.FieldValue = fmt.Sprintf("%c", c.FieldValue.(string)[0])
	if ctx := alpha(tmp); ctx != nil {
		ctx.SetTemplate("firstCharAlpha")
		return ctx
	}

	// last char can't be dash
	if strings.HasSuffix(c.FieldValue.(string), "_") {
		ctx.SetTemplate("lastUnderline")
		return ctx
	}
	return nil
}

func email(c CheckerContext) *ErrContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := ctx.FieldValue.(string)
	if value == "" {
		return nil
	}
	emailPattern := `^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`
	if !regexp.MustCompile(emailPattern).MatchString(value) {
		return ctx
	}
	return nil
}

func ipv4(c CheckerContext) *ErrContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := ctx.FieldValue.(string)
	if value == "" {
		return nil
	}
	ipv4Pattern := regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
	if !ipv4Pattern.MatchString(value) {
		return ctx
	}
	return nil
}

// MobilePattern is used to check the mobile phone.
var MobilePattern = regexp.MustCompile(`^((\+86)|(86))?(1(([35][0-9])|[8][0-9]|[7][056789]|[4][579]|99))\d{8}$`)

func mobile(c CheckerContext) *ErrContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := ctx.FieldValue.(string)
	if value == "" {
		return nil
	}
	if !MobilePattern.MatchString(value) {
		return ctx
	}
	return nil
}

func tel(c CheckerContext) *ErrContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := ctx.FieldValue.(string)
	if value == "" {
		return nil
	}
	if !regexp.MustCompile(`^(0\d{2,3}(\-)?)?\d{7,8}$`).MatchString(value) {
		return ctx
	}
	return nil
}

func phone(c CheckerContext) *ErrContext {
	telErrCtx := tel(c)
	mobileErrCtx := mobile(c)
	if telErrCtx == nil || mobileErrCtx == nil {
		return nil
	}

	if telErrCtx != nil {
		return telErrCtx
	}
	return mobileErrCtx
}

func idCard(c CheckerContext) *ErrContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.FieldValue).Kind() != reflect.String {
		return MakeValueTypeError(c)
	}

	value := ctx.FieldValue.(string)
	if value == "" {
		return nil
	}
	if !regexp.MustCompile(`(^\d{15}$)|(^\d{17}([0-9X])$)`).MatchString(value) {
		return ctx
	}
	return nil
}
