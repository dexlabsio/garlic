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

func (st *stackTrace) Key() string {
	return "stacktrace"
}

func (st *stackTrace) Value() any {
	return st.stackTrace
}

func (st *stackTrace) Visibility() Visibility {
	return RESTRICT
}

func (st *stackTrace) Insert(other Entry) Entry {
	return other
}

func (st *stackTrace) Opt() Opt {
	return func(e *ErrorT) {
		e.Insert(st)
	}
}
