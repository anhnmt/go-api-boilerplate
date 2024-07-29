package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/cmd/api-server/config"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/base"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/util"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := fx.New(
		fx.WithLogger(logger.NewFxLogger),
		fx.Provide(
			util.ProvideCtx(ctx),
		),
		fx.Invoke(
			base.GormAutoMigrate,
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
