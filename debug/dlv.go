//go:build debug
// +build debug

package debug

import (
	"fmt"
	"runtime"
)

// Breakpoint sends a SIGTRAP signal to the current process.
// When running under a debugger like Delve, this will pause execution.
// DO NOT COMMIT CODE WITH THIS FUNCTION.
func Breakpoint() {
	fmt.Println("DEBUG BREAKPOINT")
	runtime.Breakpoint()
}
