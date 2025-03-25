package validator

import (
	"reflect"
	"regexp"
	"strings"

	val "github.com/go-playground/validator/v10"
)

var SimpleValidator *val.Validate

const (
	alphaSpaceRegexString string = "^[a-zA-Z ]+$"
)

func New() *val.Validate {
	validate := val.New()

	// Using the names which have been specified for JSON representations of structs, rather than normal Go field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := validate.RegisterValidation("alpha_space", isAlphaSpace); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("is_safe_path", isSafePath); err != nil {
		panic(err)
	}

	return validate
}

func ParseValidationErrors(err error) error {
	if err == nil {
		return nil
	}

	valErrs, ok := err.(val.ValidationErrors)
	if !ok {
	}

	return NewValidationError("validation error", ValidationErrors(valErrs))
}

func isAlphaSpace(fl val.FieldLevel) bool {
	reg := regexp.MustCompile(alphaSpaceRegexString)
	return reg.MatchString(fl.Field().String())
}

func isSafePath(fl val.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), "..")
}

func init() {
	SimpleValidator = New()
}
