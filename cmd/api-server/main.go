package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/cmd/api-server/config"
	"github.com/anhnmt/go-api-boilerplate/gen/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/base"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/otel"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/permission"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/rbac"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/redis"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/util"
	"github.com/anhnmt/go-api-boilerplate/internal/server"
	"github.com/anhnmt/go-api-boilerplate/internal/server/grpc"
	"github.com/anhnmt/go-api-boilerplate/internal/server/interceptor"
	"github.com/anhnmt/go-api-boilerplate/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := fx.New(
		fx.WithLogger(logger.NewFxLogger),
		fx.Provide(
			util.ProvideCtx(ctx),
			gormgen.Use,
		),
		fx.Invoke(
			base.AutoMigrate,
		),
		config.Module,
		logger.Module,
		postgres.Module,
		redis.Module,
		otel.Module,
		rbac.Module,
		permission.Module,
		interceptor.Module,
		grpc.Module,
		service.Module,
		server.Module,
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
