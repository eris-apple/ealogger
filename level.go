package ealogger

import (
	"github.com/charmbracelet/log"
	"go.uber.org/zap/zapcore"
	"math"
)

type Level int32

const (
	DebugLevel Level = iota - 2
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	UnselectedLevel
)

func (l Level) IsEnabled(level Level) bool {
	return level >= l
}

// String returns the string representation of the level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case UnselectedLevel:
		return "unselected"
	default:
		return ""
	}
}

func (l Level) toZap() zapcore.Level {
	switch l {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	case UnselectedLevel:
		return zapcore.InfoLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l Level) toCharmbracelet() log.Level {
	switch l {
	case DebugLevel:
		return log.DebugLevel
	case InfoLevel:
		return log.InfoLevel
	case WarnLevel:
		return log.WarnLevel
	case ErrorLevel:
		return log.ErrorLevel
	case FatalLevel:
		return log.FatalLevel
	case UnselectedLevel:
		return math.MaxInt32
	default:
		return log.InfoLevel
	}
}
