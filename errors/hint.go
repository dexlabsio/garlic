package errors

import "fmt"

type hint struct {
	message string
}

func Hint(template string, args ...any) *hint {
	return &hint{
		message: fmt.Sprintf(template, args...),
	}
}

func (h *hint) Opt(e *ErrorT) {
	e.Details["hint"] = h.message
}
