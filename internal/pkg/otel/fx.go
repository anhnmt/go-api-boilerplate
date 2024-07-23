package otel

import (
	"go.uber.org/fx"
)

// Module provided to fx
var Module = fx.Module("open_telemetry", fx.Invoke(
	New,
))
