package ealogger

import (
	"encoding/json"
	"github.com/eris-apple/ealogger/ealogger/adapters"
	"github.com/eris-apple/ealogger/ealogger/shared"
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
	Infon(traceName string, args ...any)
	Infof(traceName string, args ...any)
	Debugn(traceName string, args ...any)
	Debugf(traceName string, args ...any)
	Warnn(traceName string, args ...any)
	Warnf(traceName string, args ...any)
	Errorn(traceName string, args ...any)
	Errorf(traceName string, args ...any)
	Fataln(traceName string, args ...any)
	Fatalf(traceName string, args ...any)
}

type Logger struct {
	adapters []adapters.Adapter
}

func (l *Logger) Log(log shared.Log) {
	for _, adapter := range l.adapters {
		logCopy := shared.NewLogCopy(log)
		adapter.Log(logCopy)
	}
}

func (l *Logger) WithFields(fields shared.LogField) *Entry {
	entry := NewEntry(l)
	entry.WithFields(fields)
	return entry
}

func (l *Logger) WithField(key string, value interface{}) *Entry {
	entry := NewEntry(l)
	entry.WithField(shared.LogField{key: value})
	return entry
}

func (l *Logger) WithName(traceName string) *Entry {
	entry := NewEntry(l)
	entry.WithName(traceName)
	return entry
}

func (l *Logger) WithError(err error) *Entry {
	entry := NewEntry(l)
	entry.WithError(err)
	return entry
}

func (l *Logger) Print(args ...any) {
	l.Log(shared.NewDefaultLog(shared.UnselectedLevel, args...))
}

func (l *Logger) Printf(format string, args ...any) {
	l.Log(shared.NewDefaultLogf(shared.UnselectedLevel, format, args...))
}

func (l *Logger) Info(args ...any) {
	l.Log(shared.NewDefaultLog(shared.InfoLevel, args...))
}

func (l *Logger) Infon(traceName string, args ...any) {
	l.Log(shared.NewDefaultLogn(shared.InfoLevel, traceName, args...))
}

func (l *Logger) Infof(format string, args ...any) {
	l.Log(shared.NewDefaultLogf(shared.InfoLevel, format, args...))
}

func (l *Logger) Debug(args ...any) {
	l.Log(shared.NewDefaultLog(shared.DebugLevel, args...))
}

func (l *Logger) Debugn(traceName string, args ...any) {
	l.Log(shared.NewDefaultLogn(shared.DebugLevel, traceName, args...))
}

func (l *Logger) Debugf(format string, args ...any) {
	l.Log(shared.NewDefaultLogf(shared.DebugLevel, format, args...))
}

func (e *Logger) DebugJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		e.Log(shared.NewDefaultLog(shared.DebugLevel, "error with marshaling struct"))
		return
	}

	e.Log(shared.NewDefaultLog(shared.DebugLevel, string(jsonData)))
}

func (e *Logger) DebugnJSON(traceName string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		e.Log(shared.NewDefaultLogn(shared.DebugLevel, traceName, "error with marshaling struct"))
		return
	}

	e.Log(shared.NewDefaultLogn(shared.DebugLevel, traceName, string(jsonData)))
}

func (l *Logger) Warn(args ...any) {
	l.Log(shared.NewDefaultLog(shared.WarnLevel, args...))
}

func (l *Logger) Warnn(traceName string, args ...any) {
	l.Log(shared.NewDefaultLogn(shared.WarnLevel, traceName, args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.Log(shared.NewDefaultLogf(shared.WarnLevel, format, args...))
}

func (l *Logger) Error(args ...any) {
	l.Log(shared.NewDefaultLog(shared.ErrorLevel, args...))
}

func (l *Logger) Errorn(traceName string, args ...any) {
	l.Log(shared.NewDefaultLogn(shared.ErrorLevel, traceName, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.Log(shared.NewDefaultLogf(shared.ErrorLevel, format, args...))
}

func (l *Logger) Fatal(args ...any) {
	l.Log(shared.NewDefaultLog(shared.FatalLevel, args...))
}

func (l *Logger) Fataln(traceName string, args ...any) {
	l.Log(shared.NewDefaultLogn(shared.FatalLevel, traceName, args...))
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.Log(shared.NewDefaultLogf(shared.FatalLevel, format, args...))
}

func NewLoggerWithMode(mode Mode) *Logger {
	return NewLogger(setupDefaultLogger(mode)...)
}

func NewLogger(adapters ...adapters.Adapter) *Logger {
	return &Logger{
		adapters: adapters,
	}
}

func setupDefaultLogger(mode Mode) (adp []adapters.Adapter) {
	var consoleLevel, fileLevel shared.Level

	switch mode {
	case DevMode:
		consoleLevel = shared.DebugLevel
		fileLevel = shared.DebugLevel
	case DebugMode:
		consoleLevel = shared.DebugLevel
		fileLevel = shared.InfoLevel
	case ProdMode:
		consoleLevel = shared.InfoLevel
		fileLevel = shared.WarnLevel
	}

	adp = append(adp, adapters.NewDefaultFileAdapterWithLevel(fileLevel))
	adp = append(adp, adapters.NewDefaultConsoleAdapterWithLevel(consoleLevel))

	return
}
