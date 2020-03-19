package govalid

type valid struct {
	Fields []Field
	Errors []errContext
}

// CheckFunc is the type of checker function.
type CheckFunc func(c ruleContext) *errContext

// Checkers is the function of checkers.
var Checkers map[string]CheckFunc

func init() {
	Checkers = map[string]CheckFunc{
		"required": require,
		"min":      min,
		"max":      max,
		"minlen":	minlen,
		"maxlen":	maxlen,
		//"range":        Range,
		"alpha":        alpha,
		"alphanumeric": alphaNumeric,
		"alphadash":    alphaDash,
		"username":     userName,
		"email":        email,
		"ipv4":         ipv4,
		"mobile":       mobile,
		"tel":          tel,
		"phone":        phone,
		"idcard":       idCard,
	}
}

// Check: check the struct value.
func (v *valid) Check() bool {
	for _, field := range v.Fields {
		for _, r := range field.ruleCtx {
			ruleFunc, exist := Checkers[r.rule]
			if !exist {
				// checker function not found
				v.Errors = append(v.Errors, errContext{
					Field:      field.name,
					Label:      field.label,
					Tmpl:       GetErrorTemplate("_ruleNotFound"),
					Message:    GetErrorTemplate("_ruleNotFound"),
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
				v.Errors = append(v.Errors, *err)
			}
		}
	}
	// if the error is not empty, return false
	return len(v.Errors) == 0
}
