package exampes

import (
	"github.com/eris-apple/ealogger/ealogger"
	"github.com/eris-apple/ealogger/ealogger/adapters"
	"github.com/eris-apple/ealogger/ealogger/shared"
)

func initLoggerWithMode() {
	logger := ealogger.NewLoggerWithMode(ealogger.DebugMode)
	logger.Info("Hello world!")
}

func initLoggerManually() {
	logger := ealogger.NewLogger(
		adapters.NewDefaultConsoleAdapterWithLevel(shared.DebugLevel),
		adapters.NewDefaultFileAdapterWithLevel(shared.DebugLevel),
	)

	logger.Info("Hello world!")
}

func initLoggerWithCustomConfig() {
	messageColor := "#e3e3e3"
	timestampColor := "#8a8a8a"
	consoleConfig := &adapters.ConsoleConfig{
		Enable: true,
		Level:  shared.DebugLevel,
		Colors: &adapters.ConsoleColorConfig{
			MessageColor:   &messageColor,
			TimestampColor: &timestampColor,
		},
	}

	logger := ealogger.NewLogger(
		adapters.NewConsoleAdapter(consoleConfig),
		adapters.NewDefaultFileAdapter(),
		/* defaults:
		&FileConfig{
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
		*/
	)

	logger.Info("Hello world!")
}
