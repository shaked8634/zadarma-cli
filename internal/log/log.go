package log

import (
	"fmt"
	"os"
	"sync/atomic"
)

// simple package-level logger with a debug switch.
// Use SetDebug to enable/disable debug output across the app.

var debugFlag atomic.Bool

// SetDebug enables or disables debug logging globally.
func SetDebug(enabled bool) { debugFlag.Store(enabled) }

// IsDebug reports whether debug logging is enabled.
func IsDebug() bool { return debugFlag.Load() }

// Debugf prints a formatted debug message to stderr when debug is enabled.
func Debugf(format string, args ...any) {
	if !debugFlag.Load() {
		return
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
}

// Infof prints a formatted info message to stderr.
func Infof(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[INFO] "+format+"\n", args...)
}

// Errorf prints a formatted error message to stderr.
func Errorf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[ERROR] "+format+"\n", args...)
}
