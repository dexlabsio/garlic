package errors

type UnexpectedError struct {
	*ErrorT
}

func NewUnexpectedError(message string, opts ...Opt) *UnexpectedError {
	return &UnexpectedError{
		New(message, opts...),
	}
}

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

type NotFoundError struct {
	*UserError
}

func NerNotFoundError(message string, opts ...Opt) *NotFoundError {
	return &NotFoundError{
		NewUserError(message, opts...),
	}
}
