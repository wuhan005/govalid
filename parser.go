package govalid

import (
	"reflect"
	"strings"

	"golang.org/x/text/language"
)

var (
	// RulesField is the validation rules tag's name.
	RulesField = "valid"
	// LabelField is the field label tag's name.
	LabelField = "label"
	// MessageField is the error message tag's name.
	MessageField = "msg"
)

// structField is one of the struct field.
type structField struct {
	name  string
	typ   reflect.Type
	value interface{}

	label        string
	errorMessage string

	rawRules string
	rules    []*rule
}

// parseStruct parses the given struct field.
func parseStruct(structType reflect.Type, structValue reflect.Value, languageTag language.Tag) []*structField {
	fields := make([]*structField, 0)
	rulesSets := make(map[string][]*rule)

	// Check if is a struct slice, and parse each struct.
	if structType.Kind() == reflect.Slice {
		for i := 0; i < structValue.Len(); i++ {
			structFields := parseStruct(structType.Elem(), structValue.Index(i), languageTag)
			fields = append(fields, structFields...)
		}
		return fields
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Check if the field is a struct slice.
		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
			for j := 0; j < structValue.Field(i).Len(); j++ {
				fields = append(fields, parseStruct(field.Type.Elem(), structValue.Field(i).Index(j), languageTag)...)
			}
		}

		// Check if the field is a struct.
		if field.Type.Kind() == reflect.Struct {
			fields = append(fields, parseStruct(field.Type, structValue.Field(i), languageTag)...)
		}

		// Check if this field has a validator tag.
		rawRules, ok := field.Tag.Lookup(RulesField)
		if !ok {
			continue
		}

		name := field.Name
		// Check if this field has a customized label name.
		label := name
		labelField := LabelField
		// We accept user specified language tag.
		// e.g. `label:"Name" label-en:"Name" label-zh:"姓名"`
		if languageTag.String() != "" {
			labelField += "-" + languageTag.String()
		}

		if labelValue, ok := structType.Field(i).Tag.Lookup(labelField); ok {
			label = labelValue
		} else {
			if labelValue, ok = structType.Field(i).Tag.Lookup(LabelField); ok {
				label = labelValue
			}
		}

		var errorMessage string
		if messageValue, ok := structType.Field(i).Tag.Lookup(MessageField); ok {
			errorMessage = messageValue
		}
		typ := structValue.Field(i).Type()
		value := structValue.Field(i).Interface()

		// Parse validation rules.
		// We store every field's rules in a map, so we can only parse the same rules once.
		var rules []*rule
		if rulesSet, ok := rulesSets[rawRules]; ok {
			rules = rulesSet
		} else {
			rules = parseRules(rawRules)
			rulesSets[rawRules] = rules
		}

		fields = append(fields, &structField{
			name:         name,
			typ:          typ,
			value:        value,
			label:        label,
			errorMessage: errorMessage,
			rawRules:     rawRules,
			rules:        rules,
		})
	}

	return fields
}

// rule is a single validator rule context of a struct field.
type rule struct {
	checker string
	params  []string
}

func parseRules(rawRules string) []*rule {
	rules := make([]*rule, 0)

	segments := strings.Split(rawRules, ";")
	for _, segment := range segments {
		if segment == "" {
			continue
		}

		if strings.Contains(segment, ":") {
			// key - value
			kv := strings.SplitN(segment, ":", 2)
			if kv[0] == "" {
				continue
			}
			rules = append(rules, &rule{
				checker: kv[0],
				params:  strings.Split(kv[1], ","),
			})
		} else {
			// value
			rules = append(rules, &rule{
				checker: segment,
			})
		}
	}
	return rules
}
