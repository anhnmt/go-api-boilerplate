package config

import (
	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
)

func New() (Config, error) {
	cfg := Config{}
	err := config.Load(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func appConfig(c Config) config.App {
	return c.App
}

func loggerConfig(c Config) logger.Config {
	return c.Logger
}

func postgresConfig(c Config) postgres.Config {
	return c.Postgres
}

var Module = fx.Module("config", fx.Provide(
	New,
	appConfig,
	loggerConfig,
	postgresConfig,
))
