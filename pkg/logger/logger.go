// Package logger provides structured logging functionality using Go's slog package.
// It configures a JSON-based logger with configurable log levels and source information.
package logger

import (
	"log/slog"
	"os"
)

// NewLogger initializes and configures the global logger with the specified settings.
func NewLogger(isDebug bool, addSource bool) {
	level := slog.LevelInfo
	if isDebug {
		level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
	})

	slog.SetDefault(slog.New(handler))
}
