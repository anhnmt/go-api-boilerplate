package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/anhnmt/go-api-boilerplate/cmd/api-server/config"
	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/redis"
	"github.com/anhnmt/go-api-boilerplate/internal/server"
	"github.com/anhnmt/go-api-boilerplate/internal/server/grpc"
)

var signals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(fmt.Sprintf("Failed get config: %v", err))
	}

	logger.New(cfg.Log)

	_, err = maxprocs.Set(maxprocs.Logger(log.Info().Msgf))
	if err != nil {
		panic(fmt.Sprintf("Failed set maxprocs: %v", err))
	}

	log.Info().Msg("Starting application")

	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	defer cancel()

	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		panic(fmt.Sprintf("Failed connect to database: %v", err))
	}

	rdb, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		panic(fmt.Sprintf("Failed connect to redis: %v", err))
	}

	grpcSrv, err := grpc.New(gormgen.Use(db.DB), rdb, cfg.Server.Grpc, cfg.JWT)

	server, err := server.New(grpcSrv)
	if err != nil {
		panic(fmt.Sprintf("Failed new server: %v", err))
	}

	go func() {
		err = server.Start(ctx, cfg.Server)
		if err != nil {
			panic(fmt.Sprintf("Failed start server: %v", err))
		}
	}()

	select {
	case done := <-ctx.Done():
		log.Info().Any("done", done).Msg("ctx.Done")
	}

	_ = db.Close()
	_ = rdb.Close()

	log.Info().Msg("Gracefully shutting down")
}
