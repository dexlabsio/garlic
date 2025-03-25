package validator

import (
	"fmt"

	val "github.com/go-playground/validator/v10"
	"github.com/dexlabsio/garlic/pkg/errors"
)

type validationErrors struct {
	errs val.ValidationErrors
}

func ValidationErrors(errs val.ValidationErrors) *validationErrors {
	return &validationErrors{errs}
}

func (verrs *validationErrors) Key() string {
	return "validation_errors"
}

func (verrs *validationErrors) Value() any {
	errors := make([]string, len(verrs.errs))
	for i, e := range verrs.errs {
		switch e.Tag() {
		case "required":
			errors[i] = fmt.Sprintf("%s is a required field", e.Field())
		case "max":
			errors[i] = fmt.Sprintf("%s must be a maximum of %s in length", e.Field(), e.Param())
		case "url":
			errors[i] = fmt.Sprintf("%s must be a valid URL", e.Field())
		case "alpha_space":
			errors[i] = fmt.Sprintf("%s can only contain alphabetic and space characters", e.Field())
		case "datetime":
			if e.Param() == "2006-01-02" {
				errors[i] = fmt.Sprintf("%s must be a valid date", e.Field())
			} else {
				errors[i] = fmt.Sprintf("%s must follow %s format", e.Field(), e.Param())
			}
		default:
			errors[i] = fmt.Sprintf("something wrong on %s; %s", e.Field(), e.Tag())
		}
	}

	return errors
}

func (verrs *validationErrors) Visibility() errors.Visibility {
	return errors.PUBLIC
}

func (verrs *validationErrors) Insert(other errors.Opt) errors.Opt {
	if other == nil {
		return verrs
	}

	otherVerrs, ok := other.(*validationErrors)
	if !ok {
		panic("type mismatch inserting simpleTrace opt")
	}

	verrs.errs = append(verrs.errs, otherVerrs.errs...)

	return verrs
}
