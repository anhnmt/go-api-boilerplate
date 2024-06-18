// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package service

import (
	"connectrpc.com/vanguard"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	"github.com/anhnmt/go-api-boilerplate/internal/service/user/transport/grpc"
	"net/http"
)

// Injectors from wire.go:

func New(mux *http.ServeMux, cfg config.Grpc, services *[]*vanguard.Service) error {
	userServiceHandler := usergrpc.New(services)
	error2 := initServices(mux, cfg, userServiceHandler)
	return error2
}