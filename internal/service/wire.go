//go:build wireinject
// +build wireinject

package service

import (
	"github.com/google/wire"
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	authbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/auth/business"
	authgrpc "github.com/anhnmt/go-api-boilerplate/internal/service/auth/transport/grpc"
	sessioncommand "github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/command"
	userbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/user/business"
	usercommand "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/command"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
	usergrpc "github.com/anhnmt/go-api-boilerplate/internal/service/user/transport/grpc"
)

func New(grpcSrv *grpc.Server, gormQuery *gormgen.Query, cfg config.JWT) error {
	wire.Build(
		usercommand.New,
		userquery.New,
		userbusiness.New,
		usergrpc.New,
		sessioncommand.New,
		authbusiness.New,
		authgrpc.New,
		initServices,
	)

	return nil
}
