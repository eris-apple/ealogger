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

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

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
	Fatal(v ...interface{})
	Panic(v ...interface{})
}

type Logger struct {
	DefaultLogger
	Mode          Mode
	fileLogger    *zap.Logger
	consoleLogger *log.Logger
}

func (l *Logger) Default() *Logger {
	return l
}

func (l *Logger) Println(v ...interface{}) {
	l.consoleLogger.Info(l.format(v...))
	l.fileLogger.Info(l.format(v...))
}

func (l *Logger) Debug(v ...interface{}) {
	if l.Mode == DevMode || l.Mode == DebugMode {
		l.consoleLogger.Debug(l.format(v...))
		l.fileLogger.Debug(l.format(v...))
	}
}

func (l *Logger) DebugT(trace string, v ...interface{}) {
	if l.Mode == DevMode || l.Mode == DebugMode {
		l.consoleLogger.Debug(l.formatT(false, trace, v...))
		l.fileLogger.Debug(l.formatT(true, trace, v...))
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.consoleLogger.Info(l.format(v...))
	l.fileLogger.Info(l.format(v...))
}

func (l *Logger) InfoT(trace string, v ...interface{}) {
	l.consoleLogger.Info(l.formatT(false, trace, v...))
	l.fileLogger.Info(l.formatT(true, trace, v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.consoleLogger.Warn(l.format(v...))
	l.fileLogger.Warn(l.format(v...))
}

func (l *Logger) WarnT(trace string, v ...interface{}) {
	l.consoleLogger.Warn(l.formatT(false, trace, v...))
	l.fileLogger.Warn(l.formatT(true, trace, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.consoleLogger.Error(l.format(v...))
	l.fileLogger.Error(l.format(v...))
}

func (l *Logger) ErrorT(trace string, v ...interface{}) {
	l.consoleLogger.Error(l.formatT(false, trace, v...))
	l.fileLogger.Error(l.formatT(true, trace, v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.consoleLogger.Fatal(l.format(v...))
	l.fileLogger.Fatal(l.format(v...))
}

func (l *Logger) FatalT(trace string, v ...interface{}) {
	l.consoleLogger.Fatal(l.formatT(false, trace, v...))
	l.fileLogger.Fatal(l.formatT(true, trace, v...))
}

func (l *Logger) Panic(v ...interface{}) {
	l.consoleLogger.Error(l.format(v...))
	l.fileLogger.Error(l.format(v...))
}

func (l *Logger) PanicT(trace string, v ...interface{}) {
	l.consoleLogger.Error(l.formatT(false, trace, v...))
	l.fileLogger.Panic(l.formatT(true, trace, v...))
}

func (l *Logger) format(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}

	str := fmt.Sprintf("%v", v)
	if len(str) == 0 {
		return ""
	}

	str = str[1:]
	str = str[:len(str)-1]

	return str
}

func (l *Logger) formatT(isFile bool, trace string, v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}

	str := fmt.Sprintf("%v", v)
	if len(str) == 0 {
		return ""
	}

	str = str[1:]
	str = str[:len(str)-1]

	if isFile {
		return fmt.Sprintf("%s %s", trace, str)
	}

	return fmt.Sprintf("%s%s%s: %s", Cyan, trace, Reset, str)
}

func NewDefaultLogger(mode string) *Logger {
	level := zap.InfoLevel
	if mode == DebugMode {
		level = zap.DebugLevel
	}

	pe := zap.NewProductionEncoderConfig()
	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.ISO8601TimeEncoder

	ioWriter := &lumberjack.Logger{
		Filename:   "logs/logs.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		LocalTime:  true,
		Compress:   false,
	}

	core := zapcore.NewCore(fileEncoder, zapcore.AddSync(ioWriter), level)
	fileLogger := zap.New(core)

	logLevel := log.InfoLevel
	if mode == DebugMode || mode == DevMode {
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

	return &Logger{
		Mode:          mode,
		fileLogger:    fileLogger,
		consoleLogger: logger,
	}
}
