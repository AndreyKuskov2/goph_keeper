package logger

import (
	"log/slog"
	"os"
)

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
