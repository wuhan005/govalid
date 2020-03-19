package govalid

import (
	"reflect"
	"strings"
)

var (
	// ValidField is the valid tag's name.
	ValidField = "valid"
	// LabelField is the label tag's name.
	LabelField = "label"
)

// Field is one of the form's field.
type Field struct {
	name    string
	label   string
	value   interface{}
	ruleCtx []ruleContext
}

type ruleContext struct {
	field  *Field
	rule   string
	params []string
	value  interface{}
}

// New return a govaild instance
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

		field := Field{
			name:  fieldName,
			label: fieldLabel,
			value: fieldValue,
		}
		field.ruleCtx = parseRules(validRules, &field)
		fields = append(fields, field)
	}

	return &valid{
		fields: fields,
		errors: make([]errContext, 0),
	}
}

func parseRules(rules string, field *Field) []ruleContext {
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
				field:  field,
				rule:   kv[0],
				params: strings.Split(kv[1], ","),
				value:  field.value,
			})
		} else {
			// value
			ctx = append(ctx, ruleContext{
				field:  field,
				rule:   segment,
				params: nil,
				value:  field.value,
			})
		}
	}
	return ctx
}
