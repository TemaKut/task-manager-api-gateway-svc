package logger

import "log/slog"

type leveler struct {
	lvl Level
}

func newLeveler(lvl Level) *leveler {
	return &leveler{
		lvl: lvl,
	}
}

func (l leveler) Level() slog.Level {
	switch l.lvl {
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}
