package logger

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("logger", fx.Provide(
	New,
))
