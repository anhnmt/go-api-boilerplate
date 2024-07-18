package authservice

import (
	"go.uber.org/fx"

	authbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/auth/business"
	authredis "github.com/anhnmt/go-api-boilerplate/internal/service/auth/repository/redis"
	authgrpc "github.com/anhnmt/go-api-boilerplate/internal/service/auth/transport/grpc"
)

// Module provided to fx
var Module = fx.Module("auth_service",
	fx.Provide(
		authbusiness.New,
		authredis.New,
	),
	fx.Invoke(
		authgrpc.New,
	),
)
