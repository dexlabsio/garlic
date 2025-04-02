package validator

import (
	"regexp"
	"strings"
)

var defaultExtendedValidations = []FieldValidator{
	NewValidation(
		"is_safe_path",
		func(field Field) bool {
			return !strings.Contains(field.Field().String(), "..")
		},
	),
	NewValidation(
		"alpha_space",
		func(field Field) bool {
			reg := regexp.MustCompile("^[a-zA-Z ]+$")
			return reg.MatchString(field.Field().String())
		},
	),
}
