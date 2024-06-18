package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/anhnmt/go-api-boilerplate/cmd/api-server/config"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
	"github.com/anhnmt/go-api-boilerplate/internal/server"
	"github.com/anhnmt/go-api-boilerplate/internal/service"
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
		log.Panic().Err(err).Msg("Failed set maxprocs")
	}

	log.Info().Msg("Starting application")

	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	defer cancel()

	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		log.Panic().Err(err).Msg("Failed new postgres")
	}

	mux := http.NewServeMux()

	// register service
	_ = service.New(mux, cfg.Server.Grpc)

	server := server.New(mux)

	go func() {
		err = server.Start(ctx, cfg.Server)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed start server")
		}
	}()

	select {
	case done := <-ctx.Done():
		log.Info().Any("done", done).Msg("ctx.Done")
	}

	_ = db.Close()

	log.Info().Msg("Gracefully shutting down")
}
