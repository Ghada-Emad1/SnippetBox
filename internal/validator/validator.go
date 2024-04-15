package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldsErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldsErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldsErrors == nil {
		v.FieldsErrors = make(map[string]string)
	}
	if _, exists := v.FieldsErrors[key]; !exists {
		v.FieldsErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)

	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChar(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
