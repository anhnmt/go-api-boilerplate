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
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
)

var signals = []os.Signal{os.Interrupt, syscall.SIGTERM}

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, signals...)

	select {
	case v := <-quit:
		log.Info().Any("v", v).Msg("signal.Notify")
	case done := <-ctx.Done():
		log.Info().Any("done", done).Msg("ctx.Done")
	}

	log.Info().Msg("Gracefully shutting down")
}
