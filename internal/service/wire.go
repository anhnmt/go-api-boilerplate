//go:build wireinject
// +build wireinject

package service

import (
	"github.com/google/wire"
	"google.golang.org/grpc"

	usergrpc "github.com/anhnmt/go-api-boilerplate/internal/service/user/transport/grpc"
)

func New(grpcSrv *grpc.Server) error {
	wire.Build(
		usergrpc.New,
		initServices,
	)

	return nil
}
