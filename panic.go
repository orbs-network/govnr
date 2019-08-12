package govnr

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime"
	"strings"
)

// Runs f() on the original goroutine; if it panics, logs the error and stack trace to the specified Errorer
// Very similar to GoOnce except doesn't start a new goroutine
func Recover(errorHandler Errorer, f func()) {
	tryOnce(errorHandler, f)
}

// this function is needed so that we don't return out of the goroutine when it panics
func tryOnce(errorHandler Errorer, f func()) {
	defer recoverPanics(errorHandler)
	f()
}

func recoverPanics(errorHandler Errorer) {
	if p := recover(); p != nil {
		errorHandler.Error(errors.Errorf("\npanic: %v\n\ngoroutine panicked at:\n%s\n\n", p, identifyPanic()))
	}
}

func identifyPanic() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(3, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	switch {
	case name != "":
		return fmt.Sprintf("%v:%v", name, line)
	case file != "":
		return fmt.Sprintf("%v:%v", file, line)
	}

	return fmt.Sprintf("pc:%x", pc)
}

