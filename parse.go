package govalid

import (
	"reflect"
	"strings"
)

var (
	ValidField = "valid"
	LabelField = "label"
)

type Field struct {
	name    string
	label   string
	value   interface{}
	ruleCtx []ruleContext
}

type ruleContext struct {
	name   string
	rule   string
	params []string
	value  interface{}
}

func New(inputStruct interface{}) *valid {
	structType := reflect.TypeOf(inputStruct)
	structValue := reflect.ValueOf(inputStruct)
	fields := make([]Field, 0)

	for i := 0; i < structType.NumField(); i++ {
		fieldName := structType.Field(i).Name
		fieldLabel := ""
		if label, exist := structType.Field(i).Tag.Lookup(LabelField); exist {
			fieldLabel = label
		}
		fieldValue := structValue.Field(i).Interface()
		validRules := structType.Field(i).Tag.Get(ValidField)

		fields = append(fields, Field{
			name:    fieldName,
			label:   fieldLabel,
			value:   fieldValue,
			ruleCtx: parseRules(validRules, fieldValue),
		})
	}

	return &valid{
		fields: fields,
		errors: make([]errContext, 0),
	}
}

func parseRules(rules string, val interface{}) []ruleContext {
	ctx := make([]ruleContext, 0)
	segments := strings.Split(rules, ";")
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		if strings.Contains(segment, ":") {
			// key - value
			kv := strings.SplitN(segment, ":", 2)
			if len(kv[0]) == 0 {
				continue
			}
			ctx = append(ctx, ruleContext{
				name:   segments[0],
				rule:   kv[0],
				params: strings.Split(kv[1], ","),
				value:  val,
			})
		} else {
			// value
			ctx = append(ctx, ruleContext{
				name:   segments[0],
				rule:   segment,
				params: nil,
				value:  val,
			})
		}
	}
	return ctx
}
