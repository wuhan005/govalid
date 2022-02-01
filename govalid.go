package govalid

import (
	"reflect"
)

// Check checks the struct value.
func Check(v interface{}) (errs []*ErrContext, ok bool) {
	structType := reflect.TypeOf(v)
	structValue := reflect.ValueOf(v)

	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
		structValue = structValue.Elem()
	}

	structFields := parseStruct(structType, structValue)

	for _, field := range structFields {
		for _, r := range field.rules {
			rule := r
			checkerName := rule.checker
			checkerContext := CheckerContext{
				StructValue: structValue,
				FieldName:   field.name,
				FieldType:   field.typ,
				FieldLabel:  field.label,
				FieldValue:  field.value,
				Rule:        rule,
			}

			checker, ok := Checkers[checkerName]
			if !ok {
				// Checker not found.
				errs = append(errs, MakeCheckerNotFoundError(checkerContext))
				continue
			}

			if err := checker(checkerContext); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs, len(errs) == 0
}
