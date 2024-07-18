package grpc

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("grpc", fx.Provide(
	New,
))
