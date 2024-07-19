package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/cmd/api-server/config"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
)

func provideCtx(ctx context.Context) func() context.Context {
	return func() context.Context {
		return ctx
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := fx.New(
		fx.WithLogger(logger.NewFxLogger),
		fx.Provide(
			provideCtx(ctx),
		),
		config.Module,
		logger.Module,
		postgres.Module,
		gormgen.Module,
	)

	if err := app.Start(ctx); err != nil {
		panic(fmt.Errorf("failed to start application: %w", err))
	}

	<-app.Done()

	if err := app.Stop(ctx); err != nil {
		panic(fmt.Errorf("failed to stop application: %w", err))
	}

	log.Info().Msg("Gracefully shutting down")
}
