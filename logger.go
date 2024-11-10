package ealogger

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type Level = zapcore.Level
type Mode = string

const (
	DevMode   Mode = "dev"
	DebugMode Mode = "debug"
	ProdMode  Mode = "prod"
)

type DefaultLogger interface {
	Info(v ...interface{})
	Debug(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Panic(v ...interface{})
}

type EALogger interface {
	DefaultLogger
	InfoT(trace string, v ...interface{})
	DebugT(trace string, v ...interface{})
	WarnT(trace string, v ...interface{})
	ErrorT(trace string, v ...interface{})
	FatalT(trace string, v ...interface{})
	PanicT(trace string, v ...interface{})
}

type Logger struct {
	DefaultLogger
	consoleLogger *log.Logger
	fileLogger    *zap.Logger

	c *LoggerConfig
}

type LoggerConfig struct {
	useConsole, useFile     bool
	consoleLevel, fileLevel Level
	ljLogger                *lumberjack.Logger
}

func (l *Logger) logToConsole(level Level, trace string, msg string) {
	if !l.c.useConsole {
		return
	}

	if l.c.consoleLevel.Enabled(level) {
		if trace != "" {
			trace = fmt.Sprintf("%s%s%s: ", lipgloss.Color("#36C"), trace, lipgloss.Color("#dedede"))
		}
		l.consoleLogger.Info(trace + msg)
	}
}

func (l *Logger) logToFile(level Level, trace string, msg string) {
	if !l.c.useFile {
		return
	}

	if l.c.fileLevel.Enabled(level) {
		if trace != "" {
			msg = fmt.Sprintf("%s: %s", trace, msg)
		}
		l.fileLogger.Info(msg)
	}
}

func (l *Logger) Info(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.InfoLevel, "", message)
	l.logToFile(zapcore.InfoLevel, "", message)
}

func (l *Logger) InfoT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.InfoLevel, trace, message)
	l.logToFile(zapcore.InfoLevel, trace, message)
}

func (l *Logger) Debug(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.DebugLevel, "", message)
	l.logToFile(zapcore.DebugLevel, "", message)
}

func (l *Logger) DebugT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.DebugLevel, trace, message)
	l.logToFile(zapcore.DebugLevel, trace, message)
}

func (l *Logger) Warn(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.WarnLevel, "", message)
	l.logToFile(zapcore.WarnLevel, "", message)
}

func (l *Logger) WarnT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.WarnLevel, trace, message)
	l.logToFile(zapcore.WarnLevel, trace, message)
}

func (l *Logger) Error(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.ErrorLevel, "", message)
	l.logToFile(zapcore.ErrorLevel, "", message)
}

func (l *Logger) ErrorT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.ErrorLevel, trace, message)
	l.logToFile(zapcore.ErrorLevel, trace, message)
}

func (l *Logger) Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.FatalLevel, "", message)
	l.logToFile(zapcore.FatalLevel, "", message)
}

func (l *Logger) FatalT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.FatalLevel, trace, message)
	l.logToFile(zapcore.FatalLevel, trace, message)
}

func (l *Logger) Panic(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.PanicLevel, "", message)
	l.logToFile(zapcore.PanicLevel, "", message)
}

func (l *Logger) PanicT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(zapcore.PanicLevel, trace, message)
	l.logToFile(zapcore.PanicLevel, trace, message)
}

func NewLoggerWithMode(mode Mode) *Logger {
	return NewLogger(setupDefaultConfig(mode))
}

func NewLogger(lc *LoggerConfig) *Logger {
	if lc == nil {
		lc = setupDefaultConfig(DevMode)
	}

	consoleLogger := setupConsoleLogger(lc.consoleLevel)
	fileLogger := setupFileLogger(lc.fileLevel, lc.ljLogger)

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
		consoleLevel = zapcore.DebugLevel
		fileLevel = zapcore.DebugLevel
	case DebugMode:
		consoleLevel = zapcore.DebugLevel
		fileLevel = zapcore.InfoLevel
	case ProdMode:
		consoleLevel = zapcore.InfoLevel
		fileLevel = zapcore.WarnLevel
	}

	return &LoggerConfig{
		consoleLevel: consoleLevel,
		fileLevel:    fileLevel,

		useConsole: true,
		useFile:    true,

		ljLogger: &lumberjack.Logger{
			Filename:   "logs/logs.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     28,
			LocalTime:  false,
			Compress:   false,
		},
	}
}

func setupConsoleLogger(level Level) *log.Logger {
	logLevel := log.InfoLevel
	if level <= zapcore.DebugLevel {
		logLevel = log.DebugLevel
	}

	styles := log.DefaultStyles()
	styles.Message = lipgloss.NewStyle().Foreground(lipgloss.Color("#dedede"))
	styles.Timestamp = lipgloss.NewStyle().Foreground(lipgloss.Color("#8a8a8a"))

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Level:           logLevel,
	})
	logger.SetStyles(styles)
	return logger
}

func setupFileLogger(level Level, ljLogger *lumberjack.Logger) *zap.Logger {
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(pe)

	ioWriter := ljLogger
	core := zapcore.NewCore(fileEncoder, zapcore.AddSync(ioWriter), level)
	return zap.New(core)
}
