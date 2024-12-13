package adapters

import (
	"github.com/eris-apple/ealogger/ealogger/shared"
)

type Adapter interface {
	Log(log shared.Log)
	Format(log *shared.Log)
}
