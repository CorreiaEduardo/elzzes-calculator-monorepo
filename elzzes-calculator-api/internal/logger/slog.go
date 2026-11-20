// Package logger provides the slog-backed logging adapter.
package logger

import (
	"io"
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
}

func New(w io.Writer) *Logger {
	return &Logger{logger: slog.New(slog.NewJSONHandler(w, nil))}
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
