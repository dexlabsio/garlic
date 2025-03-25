package errors

import (
	"fmt"
	"runtime"
)

type revTrace struct {
	revTrace string
}

func RevTrace(skip int) *revTrace {
	pc, file, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	thisTrace := fmt.Sprintf("[%s:%v] %s", file, line, funcName)

	return &revTrace{
		revTrace: thisTrace,
	}
}

func (rt *revTrace) Key() string {
	return "simpletrace"
}

func (rt *revTrace) Value() any {
	return rt.revTrace
}

func (rt *revTrace) Visibility() Visibility {
	return RESTRICT
}

func (rt *revTrace) Insert(other Opt) Opt {
	if other == nil {
		return rt
	}

	otherRevTrace, ok := other.(*revTrace)
	if !ok {
		panic("type mismatch inserting revTrace opt")
	}

	rt.revTrace = fmt.Sprintf("%s\n%s", rt.revTrace, otherRevTrace.revTrace)

	return rt
}
