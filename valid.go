package govalid

type valid struct {
	fields []Field
	errors []errContext
}

type CheckFunc func(c ruleContext) *errContext

var Checkers map[string]CheckFunc

func init() {
	Checkers = map[string]CheckFunc{
		"required": require,
		"min":      min,
		"max":      max,
		//"range":        Range,
		//"alpha":        Alpha,
		//"numeric":      Numeric,
		//"alphanumeric": AlphaNumeric,
		//"alphadash":    AlphaDash,
		//"username":     UserName,
		//"float":        Float,
		//"integer":      Integer,
		//"email":        Email,
		//"ipv4":         IPv4,
		//"mobile":       Mobile,
		//"mobile2":      Mobile2,
		//"tel":          Tel,
		//"phone":        Phone,
		//"idcard":       IDCard,
	}
}

func (v *valid) Check() bool {
	for _, field := range v.fields {
		for _, r := range field.ruleCtx {
			ruleFunc, exist := Checkers[r.rule]
			if !exist {
				// checker function not found
				v.errors = append(v.errors, errContext{
					Field:      field.name,
					Label:      field.label,
					Tmpl:       getErrorTemplate("_ruleNotFound"),
					Message:    getErrorTemplate("_ruleNotFound"),
					Value:      field.value,
					LimitValue: r.params,
				})
				continue
			}
			// execute the checker function
			err := ruleFunc(r)
			if err != nil {
				// set field name and value here.
				err.Field = field.name
				err.Label = field.label
				v.errors = append(v.errors, *err)
			}
		}
	}
	// if the error is not empty, return false
	return len(v.errors) == 0
}
