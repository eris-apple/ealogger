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

type Mode = string

const (
	DevMode   Mode = "dev"
	DebugMode Mode = "debug"
	ProdMode  Mode = "prod"
)

type DefaultLogger interface {
	Print(v ...interface{})
	Printf(v ...interface{})
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
	UseConsole, UseFile     bool
	ConsoleLevel, FileLevel Level
	LJLogger                *lumberjack.Logger
}

func (l *Logger) logToConsole(level Level, trace string, msg string) {
	if !l.c.UseConsole {
		return
	}

	if l.c.ConsoleLevel.IsEnabled(level) {
		if trace != "" {
			trace = fmt.Sprintf("%s%s%s: ", Cyan, trace, Reset)
		}

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
		case PanicLevel.String():
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
		case PanicLevel.String():
			l.fileLogger.Panic(trace + msg)
		case UnselectedLevel.String():
			l.fileLogger.Info(trace + msg)
		default:
			l.fileLogger.Info(trace + msg)
		}
	}
}

func (l *Logger) Print(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(UnselectedLevel, "", message)
	l.logToFile(UnselectedLevel, "", message)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.logToConsole(UnselectedLevel, "", message)
	l.logToFile(UnselectedLevel, "", message)
}

func (l *Logger) Info(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(InfoLevel, "", message)
	l.logToFile(InfoLevel, "", message)
}

func (l *Logger) InfoT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(InfoLevel, trace, message)
	l.logToFile(InfoLevel, trace, message)
}

func (l *Logger) Debug(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(DebugLevel, "", message)
	l.logToFile(DebugLevel, "", message)
}

func (l *Logger) DebugT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(DebugLevel, trace, message)
	l.logToFile(DebugLevel, trace, message)
}

func (l *Logger) Warn(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(WarnLevel, "", message)
	l.logToFile(WarnLevel, "", message)
}

func (l *Logger) WarnT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(WarnLevel, trace, message)
	l.logToFile(WarnLevel, trace, message)
}

func (l *Logger) Error(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(ErrorLevel, "", message)
	l.logToFile(ErrorLevel, "", message)
}

func (l *Logger) ErrorT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(ErrorLevel, trace, message)
	l.logToFile(ErrorLevel, trace, message)
}

func (l *Logger) Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(FatalLevel, "", message)
	l.logToFile(FatalLevel, "", message)
}

func (l *Logger) FatalT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(FatalLevel, trace, message)
	l.logToFile(FatalLevel, trace, message)
}

func (l *Logger) Panic(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(PanicLevel, "", message)
	l.logToFile(PanicLevel, "", message)
}

func (l *Logger) PanicT(trace string, v ...interface{}) {
	message := fmt.Sprint(v...)
	l.logToConsole(PanicLevel, trace, message)
	l.logToFile(PanicLevel, trace, message)
}

func (l *Logger) SetConsoleLogger(logger *log.Logger) {
	l.consoleLogger = logger
}

func (l *Logger) SetFileLogger(logger *zap.Logger) {
	l.fileLogger = logger
}

func NewLoggerWithMode(mode Mode) *Logger {
	return NewLogger(setupDefaultConfig(mode))
}

func NewLogger(lc *LoggerConfig) *Logger {
	if lc == nil {
		lc = setupDefaultConfig(DevMode)
	}

	consoleLogger := setupConsoleLogger(lc.ConsoleLevel)
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

	return &LoggerConfig{
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
}

func setupConsoleLogger(level Level) *log.Logger {
	styles := log.DefaultStyles()
	styles.Message = lipgloss.NewStyle().Foreground(lipgloss.Color("#dedede"))
	styles.Timestamp = lipgloss.NewStyle().Foreground(lipgloss.Color("#8a8a8a"))

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
