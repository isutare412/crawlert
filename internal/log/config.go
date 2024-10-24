package log

import (
	"fmt"
	"log/slog"
)

type Config struct {
	Format Format
	Level  Level
	Caller bool
}

type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

func (f Format) Validate() error {
	switch f {
	case FormatJSON, FormatText:
		return nil
	default:
		return fmt.Errorf("unknown log format '%s'", f)
	}
}

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

func (l Level) Validate() error {
	switch l {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		return nil
	default:
		return fmt.Errorf("unknown log level '%s'", l)
	}
}

func (l Level) SlogLevel() slog.Level {
	sl := slog.LevelInfo
	switch l {
	case LevelDebug:
		sl = slog.LevelDebug
	case LevelInfo:
		sl = slog.LevelInfo
	case LevelWarn:
		sl = slog.LevelWarn
	case LevelError:
		sl = slog.LevelError
	}
	return sl
}
