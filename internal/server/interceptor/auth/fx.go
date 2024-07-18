package authinterceptor

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("auth_interceptor", fx.Provide(
	New,
))
