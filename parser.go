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
func parseStruct(v interface{}) []*structField {
	structType := reflect.TypeOf(v)
	structValue := reflect.ValueOf(v)

	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
		structValue = structValue.Elem()
	}

	fields := make([]*structField, 0)
	rulesSets := make(map[string][]*rule)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Check if this field has a validator tag.
		rawRules, ok := field.Tag.Lookup(RulesField)
		if !ok {
			continue
		}
		// Parse validation rules.
		// We store every field's rules in a map, so we can only parse the same rules once.
		rulesSets[rawRules] = parseRules(rawRules)

		name := field.Name
		// Check if this field has a customized label name.
		label := name
		if labelValue, exist := structType.Field(i).Tag.Lookup(LabelField); exist {
			label = labelValue
		}
		typ := structValue.Field(i).Type()
		value := structValue.Field(i).Interface()

		fields = append(fields, &structField{
			name:     name,
			typ:      typ,
			value:    value,
			label:    label,
			rawRules: rawRules,
		})
	}

	// Set rules for each field.
	for _, field := range fields {
		field.rules = rulesSets[field.rawRules]
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
