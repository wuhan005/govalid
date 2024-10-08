package govalid

import (
	"reflect"

	"golang.org/x/text/language"
)

// Check checks the struct value.
func Check(v interface{}, lang ...language.Tag) (errs []*ErrContext, ok bool) {
	structType := reflect.TypeOf(v)
	structValue := reflect.ValueOf(v)
	validateMethod := structValue.MethodByName("Validate")

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
		field := field

		fieldErrorMessage := field.errorMessage

		for _, rule := range field.rules {
			rule := rule

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

	if validateMethod.IsValid() &&
		validateMethod.Type().NumOut() == 1 && validateMethod.Type().Out(0).Kind() == reflect.Interface {

		validateResult := validateMethod.Call(nil)[0].Interface()
		validateErr, ok := validateResult.(error)
		if ok && validateErr != nil {
			errs = append(errs, MakeUserDefinedError(validateErr.Error()))
		}
	}

	return errs, len(errs) == 0
}
