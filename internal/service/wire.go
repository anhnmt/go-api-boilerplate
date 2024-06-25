//go:build wireinject
// +build wireinject

package service

import (
	"github.com/google/wire"
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/gormgen"
	userbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/user/business"
	usercommand "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/command"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
	usergrpc "github.com/anhnmt/go-api-boilerplate/internal/service/user/transport/grpc"
)

func New(grpcSrv *grpc.Server, gormQuery *gormgen.Query) error {
	wire.Build(
		usercommand.New,
		userquery.New,
		userbusiness.New,
		usergrpc.New,
		initServices,
	)

	return nil
}
