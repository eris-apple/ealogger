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

func captureConsoleLogger(consoleLevel zapcore.Level) *log.Logger {
	buffer := new(bytes.Buffer)
	logger := log.NewWithOptions(buffer, log.Options{
		ReportTimestamp: false,
		Level:           log.Level(consoleLevel),
	})
	return logger
}

func setupObservedLogger(level zapcore.Level) (*zap.Logger, *observer.ObservedLogs) {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		level,
	)

	observed, logs := observer.New(zapcore.InfoLevel)
	logger := zap.New(zapcore.NewTee(core, observed))

	return logger, logs
}

func TestLogger_InfoMethods(t *testing.T) {
	consoleLogger := captureConsoleLogger(zapcore.InfoLevel)
	fileLogger, observedLogs := setupObservedLogger(zapcore.InfoLevel)
	logger := &Logger{
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
		c: &LoggerConfig{
			useConsole:   true,
			useFile:      true,
			consoleLevel: zapcore.InfoLevel,
			fileLevel:    zapcore.InfoLevel,
		},
	}

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
	consoleLogger := captureConsoleLogger(zapcore.DebugLevel)
	fileLogger, observedLogs := setupObservedLogger(zapcore.DebugLevel)
	logger := &Logger{
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
		c: &LoggerConfig{
			useConsole:   true,
			useFile:      true,
			consoleLevel: zapcore.DebugLevel,
			fileLevel:    zapcore.DebugLevel,
		},
	}

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
	consoleLogger := captureConsoleLogger(zapcore.ErrorLevel)
	fileLogger, observedLogs := setupObservedLogger(zapcore.ErrorLevel)
	logger := &Logger{
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
		c: &LoggerConfig{
			useConsole:   true,
			useFile:      true,
			consoleLevel: zapcore.ErrorLevel,
			fileLevel:    zapcore.ErrorLevel,
		},
	}

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

func TestLogger_FatalMethods(t *testing.T) {
	consoleLogger := captureConsoleLogger(zapcore.FatalLevel)
	fileLogger, observedLogs := setupObservedLogger(zapcore.FatalLevel)
	logger := &Logger{
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
		c: &LoggerConfig{
			useConsole:   true,
			useFile:      true,
			consoleLevel: zapcore.FatalLevel,
			fileLevel:    zapcore.FatalLevel,
		},
	}

	logger.Fatal("Test fatal message")
	if len(observedLogs.All()) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(observedLogs.All()))
	}
	if observedLogs.All()[0].Message != "Test fatal message" {
		t.Errorf("Expected 'Test fatal message', got %s", observedLogs.All()[0].Message)
	}

	traceMsg := "TRACE-FATAL"
	logger.FatalT(traceMsg, "Test fatal trace message")
	if observedLogs.All()[1].Message != fmt.Sprintf("%s: Test fatal trace message", traceMsg) {
		t.Errorf("Expected 'TRACE-FATAL: Test fatal trace message', got %s", observedLogs.All()[1].Message)
	}
}

func TestLogger_PanicMethods(t *testing.T) {
	consoleLogger := captureConsoleLogger(zapcore.PanicLevel)
	fileLogger, observedLogs := setupObservedLogger(zapcore.PanicLevel)
	logger := &Logger{
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
		c: &LoggerConfig{
			useConsole:   true,
			useFile:      true,
			consoleLevel: zapcore.PanicLevel,
			fileLevel:    zapcore.PanicLevel,
		},
	}

	logger.Panic("Test panic message")
	if len(observedLogs.All()) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(observedLogs.All()))
	}
	if observedLogs.All()[0].Message != "Test panic message" {
		t.Errorf("Expected 'Test panic message', got %s", observedLogs.All()[0].Message)
	}

	traceMsg := "TRACE-PANIC"
	logger.PanicT(traceMsg, "Test panic trace message")
	if observedLogs.All()[1].Message != fmt.Sprintf("%s: Test panic trace message", traceMsg) {
		t.Errorf("Expected 'TRACE-PANIC: Test panic trace message', got %s", observedLogs.All()[1].Message)
	}
}

func TestLoggerWithMode(t *testing.T) {
	logger := NewLoggerWithMode(DebugMode)

	if logger.c.consoleLevel != zapcore.DebugLevel {
		t.Errorf("Expected console level to be DebugLevel, got %v", logger.c.consoleLevel)
	}
	if logger.c.fileLevel != zapcore.InfoLevel {
		t.Errorf("Expected file level to be InfoLevel, got %v", logger.c.fileLevel)
	}
}
