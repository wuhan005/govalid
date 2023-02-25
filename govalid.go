package govalid

import (
	"golang.org/x/text/language"
	"reflect"
)

// Check checks the struct value.
func Check(v interface{}, lang ...language.Tag) (errs []*ErrContext, ok bool) {
	structType := reflect.TypeOf(v)
	structValue := reflect.ValueOf(v)

	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
		structValue = structValue.Elem()
	}

	templateLanguage := defaultTemplateLanguage
	if len(lang) > 0 {
		templateLanguage = lang[0]
	}

	structFields := parseStruct(structType, structValue, templateLanguage)

	for _, field := range structFields {
		fieldErrorMessage := field.errorMessage

		for _, r := range field.rules {
			rule := r
			checkerName := rule.checker
			checkerContext := CheckerContext{
				StructValue:      structValue,
				FieldName:        field.name,
				FieldType:        field.typ,
				FieldLabel:       field.label,
				FieldValue:       field.value,
				TemplateLanguage: templateLanguage,
				Rule:             rule,
			}

			checker, ok := Checkers[checkerName]
			if !ok {
				// Checker not found.
				errs = append(errs, MakeCheckerNotFoundError(checkerContext))
				continue
			}

			if err := checker(checkerContext); err != nil {
				// If the field's error message is not empty, use it.
				if fieldErrorMessage != "" {
					errs = append(errs, MakeUserDefinedError(fieldErrorMessage))
					break
				}

				errs = append(errs, err)
			}
		}
	}

	return errs, len(errs) == 0
}
