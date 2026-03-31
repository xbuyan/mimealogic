package main

import (
	"fmt"
	"time"
)

// Logger provides timestamped, togglable log output.
type Logger struct {
	enabled bool
}

// NewLogger creates a Logger that respects the enabled flag.
func NewLogger(enabled bool) *Logger {
	return &Logger{enabled: enabled}
}

// Log prints a formatted message with a timestamp prefix when logging is enabled.
func (l *Logger) Log(format string, args ...interface{}) {
	if !l.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Printf("[%s] %s\n", timestamp, message)
}
