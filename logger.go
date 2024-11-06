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

type LoggerFileConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LocalTime  bool
	Compress   bool
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
	l.fileLogger.Panic(l.format(v...))
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

func NewDefaultLogger(mode string, config *LoggerFileConfig) *Logger {
	level := zap.InfoLevel
	if mode == DebugMode {
		level = zap.DebugLevel
	}

	config = setDefaultLoggerConfig(config)

	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(pe)

	ioWriter := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		LocalTime:  config.LocalTime,
		Compress:   config.Compress,
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

func setDefaultLoggerConfig(config *LoggerFileConfig) *LoggerFileConfig {
	if config == nil {
		config = &LoggerFileConfig{}
	}

	if config.Filename == "" {
		config.Filename = "logs/logs.log"
	}
	if config.MaxSize == 0 {
		config.MaxSize = 10
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 3
	}
	if config.MaxAge == 0 {
		config.MaxAge = 28
	}

	config.LocalTime = config.LocalTime || true
	config.Compress = config.Compress || false

	return config
}
