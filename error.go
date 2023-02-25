package govalid

import (
	"fmt"
	"golang.org/x/text/language"
	"strings"
)

var _ error = (*ErrContext)(nil)

// ErrContext contains the error context.
type ErrContext struct {
	FieldName        string
	FieldLabel       string
	FieldValue       interface{}
	TemplateLanguage language.Tag

	fieldLimitValue interface{}
	errorTemplate   string
	errorMessage    string
}

func (e *ErrContext) Error() string {
	return e.errorMessage
}

var (
	FieldNamePlaceholder  = "{field}"
	FieldLimitPlaceholder = "{limit}"

	defaultTemplateLanguage = language.Chinese
)

// errorTemplateSet is the set of error templates for i18n purpose.
var errorTemplateSet = map[language.Tag]map[string]string{
	language.Chinese: errorTemplateChinese,
	language.English: errorTemplateEnglish,
}

// NewErrorContext return a error context.
func NewErrorContext(c CheckerContext) *ErrContext {
	errCtx := &ErrContext{
		FieldName:        c.FieldName,
		FieldLabel:       c.FieldLabel,
		FieldValue:       c.FieldValue,
		TemplateLanguage: c.TemplateLanguage,

		errorTemplate: getErrorTemplate(c.Rule.checker, c.TemplateLanguage),
	}
	errCtx.makeMessage()

	return errCtx
}

func (e *ErrContext) makeMessage() {
	msg := e.errorTemplate
	if strings.Contains(e.errorTemplate, FieldNamePlaceholder) || strings.Contains(e.errorTemplate, FieldLimitPlaceholder) {
		msg = strings.NewReplacer(
			FieldNamePlaceholder, e.FieldName,
			FieldLimitPlaceholder, fmt.Sprintf("%v", e.fieldLimitValue),
		).Replace(msg)
	}

	fieldLabelPrefix, limitValueSuffix := !strings.HasPrefix(e.errorTemplate, "{{"), !strings.HasSuffix(e.errorTemplate, "}}")
	if fieldLabelPrefix {
		msg = e.FieldLabel + msg
	} else {
		msg = msg[2:] // Remove the first two "{{"
	}
	if limitValueSuffix {
		if e.fieldLimitValue != nil {
			msg += fmt.Sprintf("%v", e.fieldLimitValue)
		}
	} else {
		msg = msg[:len(msg)-2] // Remove the last "}}"
	}

	e.errorMessage = msg
}

func (e *ErrContext) SetFieldLimitValue(v interface{}) {
	e.fieldLimitValue = v
	e.makeMessage()
}

func (e *ErrContext) SetTemplate(key string) {
	e.errorTemplate = getErrorTemplate(key, e.TemplateLanguage)
	e.makeMessage()
}

// getErrorTemplate return the template of the given rule name.
func getErrorTemplate(key string, templateLanguage language.Tag) string {
	errorTemplate, ok := errorTemplateSet[templateLanguage]
	if !ok {
		errorTemplate = errorTemplateSet[defaultTemplateLanguage]
	}

	if value, ok := errorTemplate[key]; ok {
		return value
	}
	return errorTemplate["_unknownErrorTemplate"]
}

func MakeUserDefinedError(msg string) *ErrContext {
	errCtx := &ErrContext{
		errorMessage: msg,
	}
	return errCtx
}

func MakeCheckerNotFoundError(c CheckerContext) *ErrContext {
	template := strings.TrimPrefix(getErrorTemplate("_checkerNotFound", c.TemplateLanguage), "~")

	errCtx := &ErrContext{
		FieldName:       c.FieldName,
		FieldLabel:      c.FieldLabel,
		FieldValue:      c.FieldValue,
		fieldLimitValue: c.Rule.params,
		errorTemplate:   template,
	}
	errCtx.makeMessage()
	return errCtx
}

func MakeCheckerParamError(c CheckerContext) *ErrContext {
	template := strings.TrimPrefix(getErrorTemplate("_paramError", c.TemplateLanguage), "~")

	errCtx := &ErrContext{
		FieldName:       c.FieldName,
		FieldLabel:      c.FieldLabel,
		FieldValue:      c.FieldValue,
		fieldLimitValue: c.Rule.params,
		errorTemplate:   template,
	}
	errCtx.makeMessage()
	return errCtx
}

func MakeValueTypeError(c CheckerContext) *ErrContext {
	template := strings.TrimPrefix(getErrorTemplate("_valueTypeError", c.TemplateLanguage), "~")

	errCtx := &ErrContext{
		FieldName:       c.FieldName,
		FieldLabel:      c.FieldLabel,
		FieldValue:      c.FieldValue,
		fieldLimitValue: c.Rule.params,
		errorTemplate:   template,
	}
	errCtx.makeMessage()
	return errCtx
}

func MakeFieldNotFoundError(c CheckerContext) *ErrContext {
	template := strings.TrimPrefix(getErrorTemplate("_fieldNotFound", c.TemplateLanguage), "~")

	errCtx := &ErrContext{
		FieldName:       c.FieldName,
		FieldLabel:      c.FieldLabel,
		FieldValue:      c.FieldValue,
		fieldLimitValue: c.Rule.params,
		errorTemplate:   template,
	}
	errCtx.makeMessage()
	return errCtx
}
