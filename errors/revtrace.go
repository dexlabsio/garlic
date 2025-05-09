package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type revTrace struct {
	revTrace string
}

func RevTrace() *revTrace {
	thisTrace := "[unknown:0] unknown"
	file, line, name, found := findCaller()
	if found {
		thisTrace = fmt.Sprintf("[%s:%v] %s", file, line, name)
	}

	return &revTrace{
		revTrace: thisTrace,
	}
}

func (rt *revTrace) Opt(e *ErrorT) {
	if e.Troubleshooting.ReverseTrace == nil {
		e.Troubleshooting.ReverseTrace = []string{}
	}

	e.Troubleshooting.ReverseTrace = append(e.Troubleshooting.ReverseTrace, rt.revTrace)
}

// findCaller iterates over the stack frames and returns the first frame
// whose package is not the same as the package of the caller of this function.
func findCaller() (file string, line int, name string, ok bool) {
	var pkg string
	var pc uintptr

	// Get the package of the function that called findCaller (skip=1)
	if pc0, _, _, ok0 := runtime.Caller(1); ok0 {
		if fn := runtime.FuncForPC(pc0); fn != nil {
			fullName := fn.Name() // Example: "github.com/username/mypackage.FunctionName"
			// Extract the package prefix by removing the function name
			if pos := strings.LastIndex(fullName, "."); pos != -1 {
				pkg = fullName[:pos]
			}
		}
	}

	// Iterate over the stack frames starting at skip=1 (or 2 to also skip findCaller itself)
	for i := 1; ; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok {
			// No more frames available, return false.
			return
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		name = fn.Name()
		// If the function's package is different from the original caller's package,
		// we've found the external caller.
		if !strings.HasPrefix(name, pkg) {
			return file, line, name, true
		}
	}
}
