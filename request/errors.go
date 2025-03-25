package request

import (
	"github.com/dexlabsio/garlic/errors"
)

type InvalidRequestError struct {
	*errors.UserError
}

func NewInvalidRequestError(message string, opts ...errors.Opt) *InvalidRequestError {
	return &InvalidRequestError{
		errors.NewUserError(message, opts...),
	}
}
