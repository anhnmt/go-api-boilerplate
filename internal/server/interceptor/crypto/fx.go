package cryptointerceptor

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("crypto_interceptor", fx.Provide(
	New,
))
