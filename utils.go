package govalid

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func require(c ruleContext) *errContext {
	ctx := NewErrorContext(c)

	if c.value == nil {
		ctx.SetMessage(ctx.Label + ctx.Tmpl)
		return ctx
	}
	// string is not ""
	if reflect.TypeOf(c.value).Kind() == reflect.String {
		if c.value.(string) == "" {
			ctx.SetMessage(c.field.label + ctx.Tmpl)
			return ctx
		} else {
			return nil
		}
	}
	return nil
}

func min(c ruleContext) *errContext {
	return minOrMax(c, "min")
}

func max(c ruleContext) *errContext {
	return minOrMax(c, "max")
}

func minOrMax(c ruleContext, flag string) *errContext {
	ctx := NewErrorContext(c)
	if len(c.params) != 1 {
		ctx.Tmpl = GetErrorTemplate("_paramError")
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	// check param
	limitStr := c.params[0]

	switch reflect.TypeOf(c.value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			ctx.Tmpl = GetErrorTemplate("_paramError")
			ctx.SetMessage(c.field.label + ctx.Tmpl)
			return ctx
		}

		value := reflect.ValueOf(c.value).Int()
		if flag == "min" {
			if value < limit {
				ctx.Message = c.field.label + fmt.Sprintf(ctx.Tmpl, limit)
				return ctx
			}
		} else {
			if value > limit {
				ctx.Message = c.field.label + fmt.Sprintf(ctx.Tmpl, limit)
				return ctx
			}
		}
		return nil
	case reflect.Float32, reflect.Float64:
		limit, err := strconv.ParseFloat(limitStr, 64)
		if err != nil {
			ctx.Tmpl = GetErrorTemplate("_paramError")
			ctx.SetMessage(c.field.label + ctx.Tmpl)
			return ctx
		}

		value, err := strconv.ParseFloat(c.params[0], 64)
		if err != nil {
			ctx.Tmpl = GetErrorTemplate("_valueTypeError")
			ctx.SetMessage(c.field.label + ctx.Tmpl)
			return ctx
		}
		if flag == "min" {
			if value < limit {
				ctx.Message = c.field.label + fmt.Sprintf(ctx.Tmpl, limit)
				return ctx
			}
		} else {
			if value > limit {
				ctx.Message = c.field.label + fmt.Sprintf(ctx.Tmpl, limit)
				return ctx
			}
		}
		return nil
	}
	return nil
}

func alpha(c ruleContext) *errContext {
	ctx := NewErrorContext(c)
	if reflect.ValueOf(c.value).Kind() != reflect.String {
		ctx.Tmpl = GetErrorTemplate("_valueTypeError")
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	for _, v := range c.value.(string) {
		if v < 'A' || (v > 'Z' && v < 'a') || v > 'z' {
			ctx.SetMessage(c.field.label + ctx.Tmpl)
			return ctx
		}
	}
	return nil
}

func alphaNumeric(c ruleContext) *errContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.value).Kind() != reflect.String {
		ctx.Tmpl = GetErrorTemplate("_valueTypeError")
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	for _, v := range c.value.(string) {
		if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
			ctx.SetMessage(c.field.label + ctx.Tmpl)
			return ctx
		}
	}
	return nil
}

func alphaDash(c ruleContext) *errContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.value).Kind() != reflect.String {
		ctx.Tmpl = GetErrorTemplate("_valueTypeError")
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	value := c.value.(string)
	if value == "" {
		return nil
	}

	if !regexp.MustCompile(`^\w+$`).MatchString(value) {
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	return nil
}

func userName(c ruleContext) *errContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.value).Kind() != reflect.String {
		ctx.Tmpl = GetErrorTemplate("_valueTypeError")
		ctx.SetMessage(ctx.Tmpl)
		return ctx
	}

	value := c.value.(string)
	if value == "" {
		return nil
	}

	// is alpha dash
	tmp := c
	tmp.rule = "alphaDash"
	if errCtx := alphaDash(tmp); errCtx != nil {
		errCtx.SetMessage(c.field.label + errCtx.Tmpl)
		return errCtx
	}

	// first char must be a alpha
	tmp = c
	tmp.value = fmt.Sprintf("%c", c.value.(string)[0])
	tmp.rule = "alpha"
	if errCtx := alpha(tmp); errCtx != nil {
		errCtx.Tmpl = GetErrorTemplate("firstCharAlpha")
		errCtx.SetMessage(c.field.label + errCtx.Tmpl)
		return errCtx
	}

	// last char can't be dash
	if strings.HasSuffix(c.value.(string), "_") {
		ctx.Tmpl = GetErrorTemplate("lastUnderline")
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	return nil
}

func email(c ruleContext) *errContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.value).Kind() != reflect.String {
		ctx.Tmpl = GetErrorTemplate("_valueTypeError")
		ctx.SetMessage(ctx.Tmpl)
		return ctx
	}
	value := ctx.Value.(string)
	if value == "" {
		return nil
	}
	emailPattern := `^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`
	if !regexp.MustCompile(emailPattern).MatchString(value) {
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	return nil
}

func ipv4(c ruleContext) *errContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.value).Kind() != reflect.String {
		ctx.Tmpl = GetErrorTemplate("_valueTypeError")
		ctx.SetMessage(ctx.Tmpl)
		return ctx
	}
	value := ctx.Value.(string)
	if value == "" {
		return nil
	}
	ipv4Pattern := regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
	if !ipv4Pattern.MatchString(value) {
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	return nil
}

// MobilePattern is used to check the mobile phone.
var MobilePattern = regexp.MustCompile(`^((\+86)|(86))?(1(([35][0-9])|[8][0-9]|[7][056789]|[4][579]|99))\d{8}$`)

func mobile(c ruleContext) *errContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.value).Kind() != reflect.String {
		ctx.Tmpl = GetErrorTemplate("_valueTypeError")
		ctx.SetMessage(ctx.Tmpl)
		return ctx
	}
	value := ctx.Value.(string)
	if value == "" {
		return nil
	}
	if !MobilePattern.MatchString(value) {
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	return nil
}

func tel(c ruleContext) *errContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.value).Kind() != reflect.String {
		ctx.Tmpl = GetErrorTemplate("_valueTypeError")
		ctx.SetMessage(ctx.Tmpl)
		return ctx
	}
	value := ctx.Value.(string)
	if value == "" {
		return nil
	}
	if !regexp.MustCompile(`^(0\d{2,3}(\-)?)?\d{7,8}$`).MatchString(value) {
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	return nil
}

func phone(c ruleContext) *errContext {
	telErr := tel(c)
	mobileErr := mobile(c)

	if telErr != nil && mobileErr != nil {
		ctx := NewErrorContext(c)
		ctx.Tmpl = GetErrorTemplate("phone")
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	return nil
}

func idCard(c ruleContext) *errContext {
	ctx := NewErrorContext(c)
	if reflect.TypeOf(c.value).Kind() != reflect.String {
		ctx.Tmpl = GetErrorTemplate("_valueTypeError")
		ctx.SetMessage(ctx.Tmpl)
		return ctx
	}
	value := ctx.Value.(string)
	if value == "" {
		return nil
	}
	if !regexp.MustCompile(`(^\d{15}$)|(^\d{17}([0-9X])$)`).MatchString(value) {
		ctx.SetMessage(c.field.label + ctx.Tmpl)
		return ctx
	}
	return nil
}
