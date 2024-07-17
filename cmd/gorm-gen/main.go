package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.uber.org/automaxprocs/maxprocs"
	"gorm.io/gen"

	"github.com/anhnmt/go-api-boilerplate/cmd/gorm-gen/config"
	"github.com/anhnmt/go-api-boilerplate/cmd/gorm-gen/generator"
	sessionentity "github.com/anhnmt/go-api-boilerplate/internal/model"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
)

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

	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		log.Panic().Err(err).Msg("Failed new postgres")
	}

	// Generate code
	g := gen.NewGenerator(gen.Config{
		OutPath: "./gen/gormgen",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(db.WithContext(ctx)) // reuse your gorm db

	// Generate basic type-safe DAO API
	g.ApplyBasic(
		sessionentity.User{},
		sessionentity.Session{},
	)

	// Generate Type Safe API with Dynamic SQL defined on Query interface
	g.ApplyInterface(func(generator.User) {}, sessionentity.User{})
	g.ApplyInterface(func(generator.Session) {}, sessionentity.Session{})

	// Generate the code
	g.Execute()

	_ = db.Close()

	log.Info().Msg("Gracefully shutting down")
}
