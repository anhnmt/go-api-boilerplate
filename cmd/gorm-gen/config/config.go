package config

import (
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
)

type Config struct {
	App      config.App      `mapstructure:"app"`
	Logger   logger.Config   `mapstructure:"log"`
	Postgres postgres.Config `mapstructure:"postgres"`
}
