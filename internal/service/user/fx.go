package userservice

import (
	"go.uber.org/fx"

	userbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/user/business"
	usercommand "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/command"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
	usergrpc "github.com/anhnmt/go-api-boilerplate/internal/service/user/transport/grpc"
)

// Module provided to fx
var Module = fx.Module("user_service",
	fx.Provide(
		userbusiness.New,
		usercommand.New,
		userquery.New,
	),
	fx.Invoke(
		usergrpc.New,
	),
)
