package ealogger

import (
	"bytes"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewDefaultLogger_WithDefaults(t *testing.T) {
	logger := NewDefaultLogger(ProdMode, nil)

	assert.Equal(t, ProdMode, logger.Mode)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.fileLogger)
	assert.NotNil(t, logger.consoleLogger)
}

func TestNewDefaultLogger_WithConfig(t *testing.T) {
	config := &LoggerFileConfig{
		Filename:   "test.log",
		MaxSize:    5,
		MaxBackups: 1,
		MaxAge:     7,
		LocalTime:  true,
		Compress:   true,
	}

	logger := NewDefaultLogger(DebugMode, config)

	assert.Equal(t, DebugMode, logger.Mode)
	assert.Equal(t, "test.log", config.Filename)
	assert.Equal(t, 5, config.MaxSize)
	assert.Equal(t, 1, config.MaxBackups)
	assert.Equal(t, 7, config.MaxAge)
	assert.True(t, config.LocalTime)
	assert.True(t, config.Compress)
}

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	consoleLogger := log.NewWithOptions(&buf, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.RFC3339,
		Level:           log.DebugLevel,
	})

	fileLogger, _ := zap.NewDevelopment()
	logger := &Logger{
		Mode:          DevMode,
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
	}

	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warn message")
	logger.Error("This is an error message")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "This is a debug message")
	assert.Contains(t, logOutput, "This is an info message")
	assert.Contains(t, logOutput, "This is a warn message")
	assert.Contains(t, logOutput, "This is an error message")
}

func TestLoggerFormatMethods(t *testing.T) {
	logger := NewDefaultLogger(DevMode, nil)
	trace := "TRACE_ID"
	msg := "test message"

	formatted := logger.format(msg)
	assert.Equal(t, "test message", formatted)

	formattedTrace := logger.formatT(false, trace, msg)
	expected := "\033[36mTRACE_ID\033[0m: test message"
	assert.Equal(t, expected, formattedTrace)
}

func TestSetDefaultLoggerConfig(t *testing.T) {
	config := LoggerFileConfig{}
	defaultConfig := setDefaultLoggerConfig(&config)

	assert.Equal(t, "logs/logs.log", defaultConfig.Filename)
	assert.Equal(t, 10, defaultConfig.MaxSize)
	assert.Equal(t, 3, defaultConfig.MaxBackups)
	assert.Equal(t, 28, defaultConfig.MaxAge)
	assert.True(t, defaultConfig.LocalTime)
	assert.False(t, defaultConfig.Compress)
}

func TestLoggerFatalAndPanic(t *testing.T) {
	var buf bytes.Buffer
	consoleLogger := log.NewWithOptions(&buf, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.RFC3339,
		Level:           log.DebugLevel,
	})
	fileLogger, _ := zap.NewDevelopment()

	logger := &Logger{
		Mode:          DevMode,
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but got none")
		}
	}()

	logger.Panic("This is a panic message")
	logOutput := buf.String()
	assert.Contains(t, logOutput, "This is a panic message")
}
