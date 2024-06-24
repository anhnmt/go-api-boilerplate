package grpc

import (
	"context"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
)

// InterceptorLogger adapts zerolog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		log := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			log.Debug().Msg(msg)
		case logging.LevelInfo:
			log.Info().Msg(msg)
		case logging.LevelWarn:
			log.Warn().Msg(msg)
		case logging.LevelError:
			log.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func New(cfg config.Grpc) *grpc.Server {
	logEvents := []logging.LoggableEvent{
		logging.StartCall,
		logging.FinishCall,
	}

	// log payload if enabled
	if cfg.LogPayload {
		logEvents = []logging.LoggableEvent{
			logging.PayloadReceived,
			logging.PayloadSent,
		}
	}

	logger := InterceptorLogger(log.Logger)
	validator, err := protovalidate.New(protovalidate.WithFailFast(true))
	if err != nil {
		panic(fmt.Errorf("failed to initialize validator: %w", err))
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(logger, logging.WithLogOnEvents(logEvents...)),
		recovery.StreamServerInterceptor(),
		protovalidate_middleware.StreamServerInterceptor(validator),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(logger, logging.WithLogOnEvents(logEvents...)),
		recovery.UnaryServerInterceptor(),
		protovalidate_middleware.UnaryServerInterceptor(validator),
	}

	// register grpc service server
	srv := grpc.NewServer(
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)

	if cfg.Reflection {
		// register grpc reflection
		reflection.Register(srv)
	}

	if cfg.HealthCheck {
		// register grpc health check
		healthcheck := health.NewServer()
		grpc_health_v1.RegisterHealthServer(srv, healthcheck)
	}

	return srv
}
