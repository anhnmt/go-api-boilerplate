package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.uber.org/automaxprocs/maxprocs"
	"gorm.io/gen"

	"github.com/anhnmt/go-api-boilerplate/cmd/api-server/config"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
	credentialentity "github.com/anhnmt/go-api-boilerplate/internal/service/credential/entity"
	deviceentity "github.com/anhnmt/go-api-boilerplate/internal/service/device/entity"
	sessionentity "github.com/anhnmt/go-api-boilerplate/internal/service/session/entity"
	userentity "github.com/anhnmt/go-api-boilerplate/internal/service/user/entity"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/query"
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
		OutPath: "./internal/infrastructure/persistence/postgresql",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(db.WithContext(ctx)) // reuse your gorm db

	// Generate basic type-safe DAO API
	g.ApplyBasic(
		userentity.User{},
		deviceentity.Device{},
		credentialentity.Credential{},
		sessionentity.Session{},
	)

	// Generate Type Safe API with Dynamic SQL defined on Query interface
	g.ApplyInterface(func(userquery.User) {}, userentity.User{})

	// Generate the code
	g.Execute()

	_ = db.Close()

	log.Info().Msg("Gracefully shutting down")
}
