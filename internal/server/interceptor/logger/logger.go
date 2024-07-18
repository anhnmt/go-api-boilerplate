package loggerinterceptor

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Logger zerolog.Logger
}

// New adapts zerolog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func New(p Params) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		log := p.Logger.With().Fields(fields).Logger()

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
