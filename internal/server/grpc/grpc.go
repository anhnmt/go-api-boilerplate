package grpc

import (
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	authinterceptor "github.com/anhnmt/go-api-boilerplate/internal/server/interceptor/auth"
	loggerinterceptor "github.com/anhnmt/go-api-boilerplate/internal/server/interceptor/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/service"
	authbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/auth/business"
)

func initServer(
	cfg config.Grpc,
	authBusiness *authbusiness.Business,
	service service.Service,
) *grpc.Server {
	logEvents := []logging.LoggableEvent{
		logging.StartCall,
		logging.FinishCall,
	}

	// log payload if enabled
	if cfg.LogPayload {
		logEvents = append(logEvents,
			logging.PayloadReceived,
			logging.PayloadSent,
		)
	}

	logger := loggerinterceptor.InterceptorLogger(log.Logger)
	validator, err := protovalidate.New(protovalidate.WithFailFast(true))
	if err != nil {
		panic(fmt.Errorf("failed to initialize validator: %w", err))
	}
	authInterceptor := authinterceptor.New(authBusiness)

	streamInterceptors := []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(logger, logging.WithLogOnEvents(logEvents...)),
		recovery.StreamServerInterceptor(),
		protovalidate_middleware.StreamServerInterceptor(validator),
		auth.StreamServerInterceptor(authInterceptor.AuthFunc()),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(logger, logging.WithLogOnEvents(logEvents...)),
		recovery.UnaryServerInterceptor(),
		protovalidate_middleware.UnaryServerInterceptor(validator),
		auth.UnaryServerInterceptor(authInterceptor.AuthFunc()),
	}

	// register grpc service server
	srv := grpc.NewServer(
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)

	// register service
	service.Register(srv)

	if cfg.Reflection {
		// register grpc reflection
		reflection.Register(srv)
	}

	if cfg.HealthCheck {
		// register grpc health check
		grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	}

	return srv
}
