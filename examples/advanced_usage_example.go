package exampes

import (
	"context"
	"github.com/eris-apple/ealogger/ealogger"
	"github.com/eris-apple/ealogger/ealogger/adapters"
	"github.com/eris-apple/ealogger/ealogger/shared"
)

type childStruct struct {
	ctx context.Context
	log *ealogger.Entry
}

func (s *childStruct) Init() {
	s.log.Debug("was init")
}

func newChildStruct(ctx context.Context, logInstance *ealogger.Logger) *childStruct {
	return &childStruct{
		ctx: ctx,
		log: logInstance.WithName("ChildStructName"),
	}
}

func initAdvancedUsage() {
	logger := ealogger.NewLogger(
		adapters.NewDefaultConsoleAdapterWithLevel(shared.DebugLevel),
		adapters.NewDefaultFileAdapterWithLevel(shared.DebugLevel),
	)

	logger.Info("Init my app")
	cs := newChildStruct(context.Background(), logger)
	cs.Init()
	logger.Info("App has been started")
}
