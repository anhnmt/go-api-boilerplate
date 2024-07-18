package redis

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("redis", fx.Provide(
	New,
))
