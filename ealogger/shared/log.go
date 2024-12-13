package shared

import (
	"fmt"
)

type LogField map[string]interface{}

type LogData struct {
	Fields    LogField
	Error     error
	TraceName string
	WithName  bool
}

type Log struct {
	Level   Level
	Message string
	Data    *LogData
}

func NewLogCopy(log Log) Log {
	fields := make(map[string]interface{})
	for k, v := range log.Data.Fields {
		fields[k] = v
	}

	return Log{
		Level:   log.Level,
		Message: log.Message,
		Data: &LogData{
			Fields:    fields,
			Error:     log.Data.Error,
			TraceName: log.Data.TraceName,
			WithName:  log.Data.WithName,
		},
	}
}

func NewDefaultLog(level Level, args ...interface{}) Log {
	return Log{
		Level:   level,
		Message: fmt.Sprint(args...),
		Data: &LogData{
			Fields: make(LogField),
		},
	}
}

func NewDefaultLogn(level Level, name string, args ...interface{}) Log {
	return Log{
		Level:   level,
		Message: fmt.Sprint(args...),
		Data: &LogData{
			TraceName: name,
			Fields:    make(LogField),
		},
	}
}

func NewDefaultLogf(level Level, format string, args ...interface{}) Log {
	return Log{
		Level:   level,
		Message: fmt.Sprintf(format, args...),
		Data: &LogData{
			Fields: make(LogField),
		},
	}
}
