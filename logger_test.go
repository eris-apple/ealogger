package ealogger

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"os"
	"testing"
)

func captureConsoleLogger(ConsoleLevel Level) *log.Logger {
	buffer := new(bytes.Buffer)
	logger := log.NewWithOptions(buffer, log.Options{
		ReportTimestamp: false,
		Level:           ConsoleLevel.toCharmbracelet(),
	})
	return logger
}

func setupObservedLogger(level Level) (*zap.Logger, *observer.ObservedLogs) {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		level.toZap(),
	)

	observed, logs := observer.New(level.toZap())
	logger := zap.New(zapcore.NewTee(core, observed))

	return logger, logs
}

func newLoggerInstance(level Level, fileLogger *zap.Logger, consoleLogger *log.Logger) *Logger {
	logger := NewLogger(&LoggerConfig{
		UseConsole:   true,
		UseFile:      true,
		ConsoleLevel: level,
		FileLevel:    level,
	})

	logger.SetConsoleLogger(consoleLogger)
	logger.SetFileLogger(fileLogger)

	return logger
}

func TestLogger_InfoMethods(t *testing.T) {
	consoleLogger := captureConsoleLogger(InfoLevel)
	fileLogger, observedLogs := setupObservedLogger(InfoLevel)
	logger := newLoggerInstance(InfoLevel, fileLogger, consoleLogger)

	logger.Info("Test info message")
	if len(observedLogs.All()) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(observedLogs.All()))
	}
	if observedLogs.All()[0].Message != "Test info message" {
		t.Errorf("Expected 'Test info message', got %s", observedLogs.All()[0].Message)
	}

	traceMsg := "TRACE-INFO"
	logger.InfoT(traceMsg, "Test info trace message")
	if observedLogs.All()[1].Message != fmt.Sprintf("%s: Test info trace message", traceMsg) {
		t.Errorf("Expected 'TRACE-INFO: Test info trace message', got %s", observedLogs.All()[1].Message)
	}
}

func TestLogger_DebugMethods(t *testing.T) {
	consoleLogger := captureConsoleLogger(DebugLevel)
	fileLogger, observedLogs := setupObservedLogger(DebugLevel)
	logger := newLoggerInstance(DebugLevel, fileLogger, consoleLogger)

	logger.Debug("Test debug message")
	if len(observedLogs.All()) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(observedLogs.All()))
	}
	if observedLogs.All()[0].Message != "Test debug message" {
		t.Errorf("Expected 'Test debug message', got %s", observedLogs.All()[0].Message)
	}

	traceMsg := "TRACE-DEBUG"
	logger.DebugT(traceMsg, "Test debug trace message")
	if observedLogs.All()[1].Message != fmt.Sprintf("%s: Test debug trace message", traceMsg) {
		t.Errorf("Expected 'TRACE-DEBUG: Test debug trace message', got %s", observedLogs.All()[1].Message)
	}
}

func TestLogger_ErrorMethods(t *testing.T) {
	consoleLogger := captureConsoleLogger(ErrorLevel)
	fileLogger, observedLogs := setupObservedLogger(ErrorLevel)
	logger := newLoggerInstance(ErrorLevel, fileLogger, consoleLogger)

	logger.Error("Test error message")
	if len(observedLogs.All()) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(observedLogs.All()))
	}
	if observedLogs.All()[0].Message != "Test error message" {
		t.Errorf("Expected 'Test error message', got %s", observedLogs.All()[0].Message)
	}

	traceMsg := "TRACE-ERROR"
	logger.ErrorT(traceMsg, "Test error trace message")
	if observedLogs.All()[1].Message != fmt.Sprintf("%s: Test error trace message", traceMsg) {
		t.Errorf("Expected 'TRACE-ERROR: Test error trace message', got %s", observedLogs.All()[1].Message)
	}
}

func TestLoggerWithMode(t *testing.T) {
	logger := NewLoggerWithMode(DebugMode)

	if logger.c.ConsoleLevel != DebugLevel {
		t.Errorf("Expected console level to be DebugLevel, got %v", logger.c.ConsoleLevel)
	}
	if logger.c.FileLevel != InfoLevel {
		t.Errorf("Expected file level to be InfoLevel, got %v", logger.c.FileLevel)
	}
}
