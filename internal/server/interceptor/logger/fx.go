package loggerinterceptor

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("logger_interceptor", fx.Provide(
	New,
))
