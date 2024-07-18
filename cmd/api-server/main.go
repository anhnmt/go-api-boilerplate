package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/cmd/api-server/config"
	"github.com/anhnmt/go-api-boilerplate/gen/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/redis"
	"github.com/anhnmt/go-api-boilerplate/internal/server"
	"github.com/anhnmt/go-api-boilerplate/internal/server/grpc"
	"github.com/anhnmt/go-api-boilerplate/internal/server/interceptor"
	"github.com/anhnmt/go-api-boilerplate/internal/service"
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
		fx.Provide(provideCtx(ctx)),
		fx.Provide(gormgen.Use),
		config.Module,
		logger.Module,
		postgres.Module,
		redis.Module,
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

	// cfg, err := config.New()
	// if err != nil {
	//     panic(fmt.Sprintf("Failed get config: %v", err))
	// }
	//
	// logger.New(cfg.Log)
	//
	// _, err = maxprocs.Set(maxprocs.Logger(log.Info().Msgf))
	// if err != nil {
	//     panic(fmt.Sprintf("Failed set maxprocs: %v", err))
	// }
	//
	// log.Info().Msg("Starting application")
	//
	// db, err := postgres.New(ctx, cfg.Postgres)
	// if err != nil {
	//     panic(fmt.Sprintf("Failed connect to database: %v", err))
	// }
	//
	// rdb, err := redis.New(ctx, cfg.Redis)
	// if err != nil {
	//     panic(fmt.Sprintf("Failed connect to redis: %v", err))
	// }
	//
	// grpcSrv, err := grpc.New(gormgen.Use(db.DB), rdb, cfg.Server.Grpc, cfg.JWT)
	//
	// server, err := server.New(grpcSrv, cfg.Crypto)
	// if err != nil {
	//     panic(fmt.Sprintf("Failed new server: %v", err))
	// }
	//
	// go func() {
	//     err = server.Start(ctx, cfg.Server)
	//     if err != nil {
	//         panic(fmt.Sprintf("Failed start server: %v", err))
	//     }
	// }()
	//
	// select {
	// case done := <-ctx.Done():
	//     log.Info().Any("done", done).Msg("ctx.Done")
	// }
	//
	// _ = db.Close()
	// _ = rdb.Close()
	//
}
