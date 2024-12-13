package exampes

import (
	"fmt"
	"github.com/eris-apple/ealogger/ealogger"
	"github.com/eris-apple/ealogger/ealogger/shared"
	"os"
	"strings"
	"time"
)

type testConfig struct {
	Enable bool
	Level  shared.Level
}

type testAdapter struct {
	writer *os.File
	cfg    *testConfig
}

func (a *testAdapter) Log(log shared.Log) {
	a.Format(&log)

	_, err := a.writer.Write([]byte(fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), log.Message)))
	if err != nil {
		return
	}
}

func (a *testAdapter) Format(log *shared.Log) {
	log.Message = strings.ToUpper(log.Message)
}

func newTestAdapter(cfg *testConfig) *testAdapter {
	return &testAdapter{
		writer: os.Stdout,
		cfg:    cfg,
	}
}

func initCustomAdapter() {
	cfg := &testConfig{
		Enable: true,
		Level:  shared.DebugLevel,
	}

	logger := ealogger.NewLogger(newTestAdapter(cfg))
	logger.Info("Hello world!")
}
