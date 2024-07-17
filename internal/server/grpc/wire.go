//go:build wireinject
// +build wireinject

package grpc

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	authinterceptor "github.com/anhnmt/go-api-boilerplate/internal/server/interceptor/auth"
	"github.com/anhnmt/go-api-boilerplate/internal/service"
	authbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/auth/business"
	authredis "github.com/anhnmt/go-api-boilerplate/internal/service/auth/repository/redis"
	authgrpc "github.com/anhnmt/go-api-boilerplate/internal/service/auth/transport/grpc"
	sessioncommand "github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/command"
	sessionquery "github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/query"
	userbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/user/business"
	usercommand "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/command"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
	usergrpc "github.com/anhnmt/go-api-boilerplate/internal/service/user/transport/grpc"
)

func New(gormQuery *gormgen.Query, rdb redis.UniversalClient, cfgGrpc config.Grpc, cfgJWT config.JWT) (*grpc.Server, error) {
	wire.Build(
		usercommand.New,
		userquery.New,
		userbusiness.New,
		usergrpc.New,
		sessioncommand.New,
		sessionquery.New,
		authredis.New,
		authbusiness.New,
		authgrpc.New,
		service.New,
		authinterceptor.New,
		initServer,
	)

	return &grpc.Server{}, nil
}
