package rbac

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("rbac", fx.Provide(
	New,
))
