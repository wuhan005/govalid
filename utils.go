package govalid

import (
	"fmt"
	"reflect"
	"strconv"
)

func require(c ruleContext) *errContext {
	ctx := &errContext{
		Tmpl:       getErrorTemplate(c.rule),
		Value:      c.value,
		LimitValue: c.params,
	}

	if c.value == nil {
		ctx.Message = ctx.Tmpl
		return ctx
	}
	// string is not ""
	if reflect.TypeOf(c.value).Kind() == reflect.String {
		if c.value.(string) == "" {
			ctx.Message = ctx.Tmpl
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
	ctx := &errContext{
		Tmpl:       getErrorTemplate(c.rule),
		Value:      c.value,
		LimitValue: c.params,
	}
	if len(c.params) != 1 {
		ctx.Tmpl = getErrorTemplate("_paramError")
		ctx.Message = ctx.Tmpl
		return ctx
	}
	// check param
	limitStr := c.params[0]

	switch reflect.TypeOf(c.value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			ctx.Tmpl = getErrorTemplate("_paramError")
			ctx.Message = ctx.Tmpl
			return ctx
		}

		value := reflect.ValueOf(c.value).Int()
		if flag == "min" {
			if value < limit {
				ctx.Message = fmt.Sprintf(ctx.Tmpl, limit)
				return ctx
			}
		} else {
			if value > limit {
				ctx.Message = fmt.Sprintf(ctx.Tmpl, limit)
				return ctx
			}
		}
		return nil
	case reflect.Float32, reflect.Float64:
		limit, err := strconv.ParseFloat(limitStr, 64)
		if err != nil {
			ctx.Tmpl = getErrorTemplate("_paramError")
			ctx.Message = ctx.Tmpl
			return ctx
		}

		value, err := strconv.ParseFloat(c.params[0], 64)
		if err != nil {
			ctx.Tmpl = getErrorTemplate("_paramError")
			ctx.Message = ctx.Tmpl
			return ctx
		}
		if flag == "min" {
			if value < limit {
				ctx.Message = fmt.Sprintf(ctx.Tmpl, limit)
				return ctx
			}
		} else {
			if value > limit {
				ctx.Message = fmt.Sprintf(ctx.Tmpl, limit)
				return ctx
			}
		}
		return nil
	}
	return nil
}
