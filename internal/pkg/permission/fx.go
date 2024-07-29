package permission

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("permission", fx.Provide(
	New,
))
