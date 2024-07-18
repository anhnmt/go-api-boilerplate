package interceptor

import (
	"go.uber.org/fx"

	authinterceptor "github.com/anhnmt/go-api-boilerplate/internal/server/interceptor/auth"
	cryptointerceptor "github.com/anhnmt/go-api-boilerplate/internal/server/interceptor/crypto"
	loggerinterceptor "github.com/anhnmt/go-api-boilerplate/internal/server/interceptor/logger"
)

// Module provided to fx
var Module = fx.Module("interceptor",
	loggerinterceptor.Module,
	authinterceptor.Module,
	cryptointerceptor.Module,
)
