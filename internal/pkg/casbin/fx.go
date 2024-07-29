package casbin

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("casbin", fx.Provide(
	New,
))
