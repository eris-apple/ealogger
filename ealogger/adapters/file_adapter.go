package adapters

import (
	"fmt"
	"github.com/eris-apple/ealogger/ealogger/shared"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type FileConfig struct {
	Enable bool

	Level    shared.Level
	LJLogger *lumberjack.Logger
}

type FileAdapter struct {
	writer *zap.Logger
	cfg    *FileConfig
}

func (a *FileAdapter) Log(log shared.Log) {
	if !a.cfg.Enable || !a.cfg.Level.IsEnabled(log.Level) {
		return
	}

	if log.Data.TraceName != "" {
		log.Data.TraceName = fmt.Sprintf("%s: ", log.Data.TraceName)
	}

	switch log.Level.String() {
	case shared.DebugLevel.String():
		a.writer.Debug(log.Data.TraceName + log.Message)
	case shared.InfoLevel.String():
		a.writer.Info(log.Data.TraceName + log.Message)
	case shared.WarnLevel.String():
		a.writer.Warn(log.Data.TraceName + log.Message)
	case shared.ErrorLevel.String():
		a.writer.Error(log.Data.TraceName + log.Message)
	case shared.FatalLevel.String():
		a.writer.Fatal(log.Data.TraceName + log.Message)
	case shared.UnselectedLevel.String():
		a.writer.Info(log.Data.TraceName + log.Message)
	default:
		a.writer.Info(log.Data.TraceName + log.Message)
	}
}

func (a *FileAdapter) Format(log *shared.Log) {

}

func NewFileAdapter(cfg *FileConfig) *FileAdapter {
	return &FileAdapter{
		cfg:    cfg,
		writer: newFileLogger(cfg),
	}
}

func NewDefaultFileAdapter() *FileAdapter {
	cfg := defaultFileConfig()

	return &FileAdapter{
		cfg:    cfg,
		writer: newFileLogger(cfg),
	}
}

func NewDefaultFileAdapterWithLevel(level shared.Level) *FileAdapter {
	cfg := defaultFileConfig()
	cfg.Level = level

	return &FileAdapter{
		cfg:    cfg,
		writer: newFileLogger(cfg),
	}
}

func defaultFileConfig() *FileConfig {
	return &FileConfig{
		Enable: true,
		Level:  shared.DebugLevel,
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

func newFileLogger(cfg *FileConfig) *zap.Logger {
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(pe)

	ioWriter := cfg.LJLogger
	core := zapcore.NewCore(fileEncoder, zapcore.AddSync(ioWriter), cfg.Level.ToZap())
	return zap.New(core)
}
