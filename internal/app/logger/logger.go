package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger struct {
	baseLogger *slog.Logger
}

func NewLogger(lvl Level) *Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: newLeveler(lvl),
	})

	return &Logger{baseLogger: slog.New(handler)}
}

func (l *Logger) Debugf(format string, args ...any) {
	l.baseLogger.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...any) {
	l.baseLogger.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.baseLogger.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.baseLogger.Error(fmt.Sprintf(format, args...))
}
