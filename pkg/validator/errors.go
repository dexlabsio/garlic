package validator

import (
	"github.com/dexlabsio/garlic/pkg/errors"
)

type ValidationError struct {
	*errors.UserError
}

func NewValidationError(message string, opts ...errors.Opt) *ValidationError {
	return &ValidationError{
		errors.NewUserError(message, opts...),
	}
}
