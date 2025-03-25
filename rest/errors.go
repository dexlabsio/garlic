package rest

import "github.com/dexlabsio/garlic/errors"

type NotFoundError struct {
	*errors.UserError
}

func NewNotFoundError(message string, opts ...errors.Opt) *NotFoundError {
	return &NotFoundError{
		errors.NewUserError(message, opts...),
	}
}
