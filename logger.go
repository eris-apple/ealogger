package ealogger

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
)

type Mode = string

const (
	DevMode   Mode = "dev"
	DebugMode Mode = "debug"
	ProdMode  Mode = "prod"
)

type DefaultLogger interface {
	Print(args ...any)
	Printf(args ...any)
	Info(args ...any)
	Debug(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
}

type EALogger interface {
	DefaultLogger
	InfoT(trace string, args ...any)
	DebugT(trace string, args ...any)
	WarnT(trace string, args ...any)
	ErrorT(trace string, args ...any)
	FatalT(trace string, args ...any)
}

type Field map[string]interface{}

type Logger struct {
	DefaultLogger
	consoleLogger *log.Logger
	fileLogger    *zap.Logger

	c *LoggerConfig
}
type LoggerConfig struct {
	UseConsole, UseFile     bool
	ConsoleLevel, FileLevel Level
	LJLogger                *lumberjack.Logger

	TimestampColor *string
	MessageColor   *string
	ErrorColor     *string

	LevelColors map[Level]string
}

func (l *Logger) SetConsoleLogger(logger *log.Logger) {
	l.consoleLogger = logger
}

func (l *Logger) SetFileLogger(logger *zap.Logger) {
	l.fileLogger = logger
}

func (l *Logger) logToConsole(level Level, trace string, msg string) {
	if !l.c.UseConsole {
		return
	}

	if l.c.ConsoleLevel.IsEnabled(level) {
		if trace != "" {
			trace = lipgloss.
				NewStyle().
				SetString(fmt.Sprintf("[%s]: ", trace)).
				Foreground(lipgloss.Color(l.c.LevelColors[level])).
				String()
		}

		msg = lipgloss.
			NewStyle().
			SetString(fmt.Sprintf("%s", msg)).
			Foreground(lipgloss.Color(*l.c.MessageColor)).
			String()

		switch level.String() {
		case DebugLevel.String():
			l.consoleLogger.Debug(trace + msg)
		case InfoLevel.String():
			l.consoleLogger.Info(trace + msg)
		case WarnLevel.String():
			l.consoleLogger.Warn(trace + msg)
		case ErrorLevel.String():
			l.consoleLogger.Error(trace + msg)
		case FatalLevel.String():
			l.consoleLogger.Fatal(trace + msg)
		case UnselectedLevel.String():
			l.consoleLogger.Print(trace + msg)
		default:
			l.consoleLogger.Info(trace + msg)
		}
	}
}

func (l *Logger) logToFile(level Level, trace string, msg string) {
	if !l.c.UseFile {
		return
	}

	if l.c.FileLevel.IsEnabled(level) {
		if trace != "" {
			trace = fmt.Sprintf("%s: ", trace)
		}

		switch level.String() {
		case DebugLevel.String():
			l.fileLogger.Debug(trace + msg)
		case InfoLevel.String():
			l.fileLogger.Info(trace + msg)
		case WarnLevel.String():
			l.fileLogger.Warn(trace + msg)
		case ErrorLevel.String():
			l.fileLogger.Error(trace + msg)
		case FatalLevel.String():
			l.fileLogger.Fatal(trace + msg)
		case UnselectedLevel.String():
			l.fileLogger.Info(trace + msg)
		default:
			l.fileLogger.Info(trace + msg)
		}
	}
}

func (l *Logger) Log(level Level, trace string, args ...any) {
	message := fmt.Sprint(args...)
	l.logToConsole(level, trace, message)
	l.logToFile(level, trace, message)
}

func (l *Logger) Logf(level Level, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.logToConsole(level, "", message)
	l.logToFile(level, "", message)
}

func (l *Logger) WithFields(fields Field) *Entry {
	entry := NewEntry(l)
	entry.WithFields(fields)
	return entry
}

func (l *Logger) WithField(key string, value interface{}) *Entry {
	entry := NewEntry(l)
	entry.WithField(Field{key: value})
	return entry
}

func (l *Logger) WithTrace(trace string) *Entry {
	entry := NewEntry(l)
	entry.WithTrace(trace)
	return entry
}

func (l *Logger) WithError(err error) *Entry {
	entry := NewEntry(l)
	entry.WithError(err)
	return entry
}

func (l *Logger) Print(args ...any) {
	l.Log(UnselectedLevel, "", args...)
}

func (l *Logger) Printf(format string, args ...any) {
	l.Logf(UnselectedLevel, format, args...)
}

func (l *Logger) Info(args ...any) {
	l.Log(InfoLevel, "", args...)
}

func (l *Logger) InfoT(trace string, args ...any) {
	l.Log(InfoLevel, trace, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.Logf(InfoLevel, format, args...)
}

func (l *Logger) Debug(args ...any) {
	l.Log(DebugLevel, "", args...)
}

func (l *Logger) DebugT(trace string, args ...any) {
	l.Log(DebugLevel, trace, args...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.Logf(DebugLevel, format, args...)
}

func (l *Logger) Warn(args ...any) {
	l.Log(WarnLevel, "", args...)
}

func (l *Logger) WarnT(trace string, args ...any) {
	l.Log(WarnLevel, trace, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.Logf(WarnLevel, format, args...)
}

func (l *Logger) Error(args ...any) {
	l.Log(ErrorLevel, "", args...)
}

func (l *Logger) ErrorT(trace string, args ...any) {
	l.Log(ErrorLevel, trace, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.Logf(ErrorLevel, format, args...)
}

func (l *Logger) Fatal(args ...any) {
	l.Log(FatalLevel, "", args...)
}

func (l *Logger) FatalT(trace string, args ...any) {
	l.Log(FatalLevel, trace, args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.Logf(FatalLevel, format, args...)
}

func NewLoggerWithMode(mode Mode) *Logger {
	return NewLogger(setupDefaultConfig(mode))
}

func NewLogger(lc *LoggerConfig) *Logger {
	if lc == nil {
		lc = setupDefaultConfig(DevMode)
	}

	setupDefaultLoggerColors(lc)
	consoleLogger := setupConsoleLogger(lc.ConsoleLevel, lc)
	fileLogger := setupFileLogger(lc.FileLevel, lc.LJLogger)

	return &Logger{
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
		c:             lc,
	}
}

func setupDefaultConfig(mode Mode) *LoggerConfig {
	var consoleLevel, fileLevel Level

	switch mode {
	case DevMode:
		consoleLevel = DebugLevel
		fileLevel = DebugLevel
	case DebugMode:
		consoleLevel = DebugLevel
		fileLevel = InfoLevel
	case ProdMode:
		consoleLevel = InfoLevel
		fileLevel = WarnLevel
	}

	lc := &LoggerConfig{
		ConsoleLevel: consoleLevel,
		FileLevel:    fileLevel,

		UseConsole: true,
		UseFile:    true,

		LJLogger: &lumberjack.Logger{
			Filename:   "logs/logs.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     28,
			LocalTime:  false,
			Compress:   false,
		},
	}

	setupDefaultLoggerColors(lc)
	return lc
}

func setupDefaultLoggerColors(lc *LoggerConfig) {
	timestampColor := "#8a8a8a"
	messageColor := "#e3e3e3"
	errorColor := "#af0000"

	if lc.TimestampColor == nil {
		lc.TimestampColor = &timestampColor
	}
	if lc.MessageColor == nil {
		lc.MessageColor = &messageColor
	}
	if lc.ErrorColor == nil {
		lc.ErrorColor = &errorColor
	}

	if lc.LevelColors == nil {
		lc.LevelColors = make(map[Level]string)
	}

	if lc.LevelColors[InfoLevel] == "" {
		lc.LevelColors[InfoLevel] = "#afd7ff"
	}

	if lc.LevelColors[DebugLevel] == "" {
		lc.LevelColors[DebugLevel] = "#969696"
	}

	if lc.LevelColors[WarnLevel] == "" {
		lc.LevelColors[WarnLevel] = "#ffff18"
	}

	if lc.LevelColors[ErrorLevel] == "" {
		lc.LevelColors[ErrorLevel] = "#af0000"
	}

	if lc.LevelColors[FatalLevel] == "" {
		lc.LevelColors[FatalLevel] = "#ff0000"
	}

}

func setupConsoleLogger(level Level, lc *LoggerConfig) *log.Logger {
	styles := log.DefaultStyles()
	styles.Message = lipgloss.NewStyle().Foreground(lipgloss.Color(*lc.MessageColor))
	styles.Timestamp = lipgloss.NewStyle().Foreground(lipgloss.Color(*lc.TimestampColor))
	for key, value := range lc.LevelColors {
		styles.Levels[key.toCharmbracelet()] = lipgloss.NewStyle().SetString(strings.ToUpper(key.String())).Foreground(lipgloss.Color(value))
	}

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Level:           level.toCharmbracelet(),
	})
	logger.SetStyles(styles)
	return logger
}

func setupFileLogger(level Level, LJLogger *lumberjack.Logger) *zap.Logger {
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(pe)

	ioWriter := LJLogger
	core := zapcore.NewCore(fileEncoder, zapcore.AddSync(ioWriter), level.toZap())
	return zap.New(core)
}
