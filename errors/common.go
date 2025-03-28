package errors

type UserError struct {
	*ErrorT
}

func NewUserError(message string, opts ...Opt) *UserError {
	return &UserError{
		New(message, opts...),
	}
}

type SystemError struct {
	*ErrorT
}

func NewSystemError(message string, opts ...Opt) *SystemError {
	return &SystemError{
		New(message, opts...),
	}
}
