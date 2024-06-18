//go:build wireinject
// +build wireinject

package service

import (
	"net/http"

	"connectrpc.com/vanguard"
	"github.com/google/wire"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	usergrpc "github.com/anhnmt/go-api-boilerplate/internal/service/user/transport/grpc"
)

func New(mux *http.ServeMux, cfg config.Grpc, services *[]*vanguard.Service) error {
	wire.Build(
		usergrpc.New,
		initServices,
	)

	return nil
}
