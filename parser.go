package govalid

import (
	"reflect"
	"strings"
)

var (
	// RulesField is the validation rules tag's name.
	RulesField = "valid"
	// LabelField is the field label tag's name.
	LabelField = "label"
)

// structField is one of the struct field.
type structField struct {
	name  string
	typ   reflect.Type
	value interface{}

	label string

	rawRules string
	rules    []*rule
}

// parseStruct parses the given struct field.
func parseStruct(structType reflect.Type, structValue reflect.Value) []*structField {
	fields := make([]*structField, 0)
	rulesSets := make(map[string][]*rule)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Check if the field is a struct slice.
		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
			for j := 0; j < structValue.Field(i).Len(); j++ {
				fields = append(fields, parseStruct(field.Type.Elem(), structValue.Field(i).Index(j))...)
			}
		}

		// Check if this field has a validator tag.
		rawRules, ok := field.Tag.Lookup(RulesField)
		if !ok {
			continue
		}

		name := field.Name
		// Check if this field has a customized label name.
		label := name
		if labelValue, exist := structType.Field(i).Tag.Lookup(LabelField); exist {
			label = labelValue
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
			name:     name,
			typ:      typ,
			value:    value,
			label:    label,
			rawRules: rawRules,
			rules:    rules,
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
