package validator

import (
	"fmt"

	"github.com/dexlabsio/garlic/errors"
	val "github.com/go-playground/validator/v10"
)

type validationErrors struct {
	errs val.ValidationErrors
}

func ValidationErrors(errs val.ValidationErrors) *validationErrors {
	return &validationErrors{errs}
}

func (verrs *validationErrors) Value() map[string]string {
	errors := make(map[string]string, len(verrs.errs))
	for _, e := range verrs.errs {
		field := e.Field()
		msg := fmt.Sprintf("something wrong on %s; %s", field, e.Tag())
		switch e.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is a required field", field)
		case "max":
			msg = fmt.Sprintf("%s must be a maximum of %s in length", field, e.Param())
		case "url":
			msg = fmt.Sprintf("%s must be a valid URL", field)
		case "alpha_space":
			msg = fmt.Sprintf("%s can only contain alphabetic and space characters", field)
		case "datetime":
			if e.Param() == "2006-01-02" {
				msg = fmt.Sprintf("%s must be a valid date", field)
			} else {
				msg = fmt.Sprintf("%s must follow %s format", field, e.Param())
			}
		}

		errors[field] = msg
	}

	return errors
}

func (verrs *validationErrors) Opt(e *errors.ErrorT) {
	errors := verrs.Value()

	validationDetails, ok := e.Details["validation"].(map[string]string)
	if !ok {
		validationDetails = make(map[string]string, len(errors))
		e.Details["validation"] = validationDetails
	}

	for k, v := range errors {
		validationDetails[k] = v
	}
}
