package gormgen

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("gormgen", fx.Invoke(
	New,
))
