package ealogger

import (
	"encoding/json"
	"github.com/eris-apple/ealogger/ealogger/shared"
)

type Entry struct {
	l    *Logger
	data *shared.LogData
}

func (e *Entry) Log(log shared.Log) {
	log.Data = e.data
	e.l.Log(log)

	defer func() {
		e.data.Error = nil
		e.data.Fields = nil
	}()
}

func (e *Entry) WithField(field shared.LogField) *Entry {
	e.data.Fields = field

	return e
}

func (e *Entry) WithFields(Fields shared.LogField) *Entry {
	e.data.Fields = Fields

	return e
}

func (e *Entry) WithName(traceName string) *Entry {
	e.data.TraceName = traceName
	e.data.WithName = true

	return e
}

func (e *Entry) ClearName() *Entry {
	e.data.TraceName = ""
	e.data.WithName = false

	return e
}

func (e *Entry) WithError(err error) *Entry {
	e.data.Error = err

	return e
}

func (e *Entry) Print(args ...any) {
	e.Log(shared.NewDefaultLog(shared.UnselectedLevel, args...))
}

func (e *Entry) Printf(format string, args ...any) {
	e.Log(shared.NewDefaultLogf(shared.UnselectedLevel, format, args...))
}

func (e *Entry) Info(args ...any) {
	e.Log(shared.NewDefaultLog(shared.InfoLevel, args...))
}

func (e *Entry) Infof(format string, args ...any) {
	e.Log(shared.NewDefaultLogf(shared.InfoLevel, format, args...))
}

func (e *Entry) Debug(args ...any) {
	e.Log(shared.NewDefaultLog(shared.DebugLevel, args...))
}

func (e *Entry) Debugf(format string, args ...any) {
	e.Log(shared.NewDefaultLogf(shared.DebugLevel, format, args...))
}

func (e *Entry) DebugJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		e.Log(shared.NewDefaultLog(shared.DebugLevel, "error with marshaling struct"))
		return
	}

	e.Log(shared.NewDefaultLog(shared.DebugLevel, string(jsonData)))
}

func (e *Entry) Warn(args ...any) {
	e.Log(shared.NewDefaultLog(shared.WarnLevel, args...))
}

func (e *Entry) Warnf(format string, args ...any) {
	e.Log(shared.NewDefaultLogf(shared.WarnLevel, format, args...))
}

func (e *Entry) Error(args ...any) {
	e.Log(shared.NewDefaultLog(shared.ErrorLevel, args...))
}

func (e *Entry) Errorf(format string, args ...any) {
	e.Log(shared.NewDefaultLogf(shared.ErrorLevel, format, args...))
}

func (e *Entry) Fatal(args ...any) {
	e.Log(shared.NewDefaultLog(shared.FatalLevel, args...))
}

func (e *Entry) Fatalf(format string, args ...any) {
	e.Log(shared.NewDefaultLogf(shared.FatalLevel, format, args...))
}

func NewEntry(l *Logger) *Entry {
	return &Entry{
		l: l,
		data: &shared.LogData{
			Fields: make(shared.LogField),
		},
	}
}
