package errors

import (
	"runtime/debug"
)

type stackTrace struct {
	stackTrace string
}

func StackTrace() *stackTrace {
	return &stackTrace{
		stackTrace: string(debug.Stack()),
	}
}

func (st *stackTrace) Opt(e *ErrorT) {
	e.Troubleshooting.StackTrace = st.stackTrace
}
