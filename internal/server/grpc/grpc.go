package grpc

import (
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
)

type Params struct {
	fx.In

	Config config.Grpc
	// AuthInterceptor   authinterceptor.AuthInterceptor
	LoggerInterceptor logging.Logger
}

func New(p Params) *grpc.Server {
	logEvents := []logging.LoggableEvent{
		logging.StartCall,
		logging.FinishCall,
	}

	// log payload if enabled
	if p.Config.LogPayload {
		logEvents = append(logEvents,
			logging.PayloadReceived,
			logging.PayloadSent,
		)
	}

	validator, err := protovalidate.New(protovalidate.WithFailFast(true))
	if err != nil {
		panic(fmt.Errorf("failed to initialize validator: %w", err))
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(p.LoggerInterceptor, logging.WithLogOnEvents(logEvents...)),
		recovery.StreamServerInterceptor(),
		protovalidate_middleware.StreamServerInterceptor(validator),
		// auth.StreamServerInterceptor(p.AuthInterceptor.AuthFunc()),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(p.LoggerInterceptor, logging.WithLogOnEvents(logEvents...)),
		recovery.UnaryServerInterceptor(),
		protovalidate_middleware.UnaryServerInterceptor(validator),
		// auth.UnaryServerInterceptor(p.AuthInterceptor.AuthFunc()),
	}

	// register grpc service server
	srv := grpc.NewServer(
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)

	// register service
	// p.Service.Register(srv)

	if p.Config.Reflection {
		// register grpc reflection
		reflection.Register(srv)
	}

	if p.Config.HealthCheck {
		// register grpc health check
		grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	}

	return srv
}
