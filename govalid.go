package govalid

// Check checks the struct value.
func Check(v interface{}) (errs []*ErrContext, ok bool) {
	structFields := parseStruct(v)

	for _, field := range structFields {
		for _, r := range field.rules {
			rule := r
			checkerName := rule.checker
			checkerContext := CheckerContext{
				FieldName:  field.name,
				FieldType:  field.typ,
				FieldLabel: field.label,
				FieldValue: field.value,
				Rule:       rule,
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
