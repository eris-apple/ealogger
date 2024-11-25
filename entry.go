package ealogger

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

type Entry struct {
	l *Logger

	fields    Field
	err       error
	trace     string
	withTrace bool
}

func (e *Entry) WithField(field Field) *Entry {
	e.fields = field

	return e
}

func (e *Entry) WithFields(fields Field) *Entry {
	e.fields = fields

	return e
}

func (e *Entry) WithTrace(trace string) *Entry {
	e.trace = trace
	e.withTrace = true

	return e
}

func (e *Entry) ClearTrace() *Entry {
	e.trace = ""
	e.withTrace = false

	return e
}

func (e *Entry) WithError(err error) *Entry {
	e.err = err

	return e
}

func (e *Entry) Log(level Level, trace string, args ...any) {
	message := fmt.Sprint(args...)

	if e.err != nil {
		errf := lipgloss.
			NewStyle().
			SetString(fmt.Sprintf("err=%s", e.err)).
			Foreground(lipgloss.Color(*e.l.c.ErrorColor)).
			String()
		message = fmt.Sprintf("%s %s", message, errf)
		defer func() {
			e.err = nil
		}()
	}

	if len(e.fields) > 0 {
		fieldsf := ""
		for key, field := range e.fields {
			fieldf := lipgloss.
				NewStyle().
				SetString(fmt.Sprintf("%s=%s ", key, field)).
				Foreground(lipgloss.Color(e.l.c.LevelColors[level])).
				String()
			fieldsf += fieldf
		}

		message = fmt.Sprintf("%s %s", message, fieldsf)

		defer func() {
			e.fields = nil
		}()
	}

	e.l.Log(level, trace, message)
}

func (e *Entry) Logf(level Level, trace string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)

	e.l.logToConsole(level, trace, message)
	e.l.logToFile(level, trace, message)
}

func (e *Entry) Print(args ...any) {
	e.Log(UnselectedLevel, e.trace, args...)
}

func (e *Entry) Printf(format string, args ...any) {
	e.Logf(UnselectedLevel, e.trace, format, args...)
}

func (e *Entry) Info(args ...any) {
	e.Log(InfoLevel, e.trace, args...)
}

func (e *Entry) Infof(format string, args ...any) {
	e.Logf(InfoLevel, e.trace, format, args...)
}

func (e *Entry) Debug(args ...any) {
	e.Log(DebugLevel, e.trace, args...)
}

func (e *Entry) Debugf(format string, args ...any) {
	e.Logf(DebugLevel, e.trace, format, args...)
}

func (e *Entry) Warn(args ...any) {
	e.Log(WarnLevel, e.trace, args...)
}

func (e *Entry) Warnf(format string, args ...any) {
	e.Logf(WarnLevel, e.trace, format, args...)
}

func (e *Entry) Error(args ...any) {
	e.Log(ErrorLevel, e.trace, args...)
}

func (e *Entry) Errorf(format string, args ...any) {
	e.Logf(ErrorLevel, e.trace, format, args...)
}

func (e *Entry) Fatal(args ...any) {
	e.Log(FatalLevel, e.trace, args...)
}

func (e *Entry) Fatalf(format string, args ...any) {
	e.Logf(FatalLevel, e.trace, format, args...)
}

func NewEntry(l *Logger) *Entry {
	return &Entry{
		l: l,
	}
}
