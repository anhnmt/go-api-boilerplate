package service

import (
	"net/http"

	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	"github.com/anhnmt/go-api-boilerplate/proto/pb/pbconnect"
)

var serviceNames = []string{
	pbconnect.UserServiceName,
}

func initServices(
	mux *http.ServeMux,
	cfg config.Grpc,

	_ pbconnect.UserServiceHandler,
) error {
	if cfg.Reflection {
		grpcReflect(mux)
	}

	if cfg.HealthCheck {
		grpcHealth(mux)
	}

	return nil
}

func grpcReflect(mux *http.ServeMux) {
	reflector := grpcreflect.NewStaticReflector(serviceNames...)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	// Many tools still expect the older version of the server reflection API, so
	// most servers should mount both handlers.
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}

func grpcHealth(mux *http.ServeMux) {
	checker := grpchealth.NewStaticChecker(serviceNames...)
	mux.Handle(grpchealth.NewHandler(checker))
}
