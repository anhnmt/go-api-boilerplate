package server

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("server", fx.Invoke(
	New,
))
