package adapters

import (
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/eris-apple/ealogger/ealogger/shared"
	"time"
)

type GraylogConfig struct {
	Enable bool

	Addr string
	Host string

	Level shared.Level
}

type GraylogAdapter struct {
	writer *gelf.Writer
	cfg    *GraylogConfig
}

func (a *GraylogAdapter) Log(log shared.Log) {
	if !a.cfg.Enable || !a.cfg.Level.IsEnabled(log.Level) {
		return
	}

	a.Format(&log)

	if err := a.writer.WriteMessage(&gelf.Message{
		Level:    log.Level.ToGraylog(),
		Full:     log.Data.TraceName + log.Message,
		Short:    log.Data.TraceName + log.Message,
		Host:     a.cfg.Host,
		TimeUnix: float64(time.Now().Unix()),
		Extra:    log.Data.Fields,
	}); err != nil {
		return
	}
}

func (a *GraylogAdapter) Format(log *shared.Log) {
	if log.Data != nil {
		log.Data.Fields["traceName"] = log.Data.TraceName

		if log.Data.Fields == nil {
			log.Data.Fields = map[string]interface{}{}
		}

		if log.Data.Error != nil {
			log.Data.Fields["error"] = log.Data.Error
		}
	}

}

func NewGraylogAdapter(cfg *GraylogConfig) *GraylogAdapter {
	return &GraylogAdapter{
		cfg:    cfg,
		writer: newGraylogLogger(cfg),
	}
}

func NewDefaultGraylogAdapter() *GraylogAdapter {
	cfg := defaultGraylogConfig()

	return &GraylogAdapter{
		cfg:    cfg,
		writer: newGraylogLogger(cfg),
	}
}

func NewDefaultGraylogAdapterWithLevel(level shared.Level) *GraylogAdapter {
	cfg := defaultGraylogConfig()
	cfg.Level = level

	return &GraylogAdapter{
		cfg:    cfg,
		writer: newGraylogLogger(cfg),
	}
}

func defaultGraylogConfig() *GraylogConfig {
	return &GraylogConfig{
		Enable: true,
		Level:  shared.DebugLevel,
		Addr:   "localhost:12201",
		Host:   "APP",
	}
}

func newGraylogLogger(cfg *GraylogConfig) *gelf.Writer {
	gelfWriter, err := gelf.NewWriter(cfg.Addr)
	if err != nil {
		fmt.Println("WARN: error with init graylog: ", err)
		return nil
	}

	return gelfWriter
}
