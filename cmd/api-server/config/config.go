package config

import (
	"fmt"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/redis"
)

type Config struct {
	Log      logger.Config   `mapstructure:"log"`
	Postgres postgres.Config `mapstructure:"postgres"`
	Redis    redis.Config    `mapstructure:"redis"`
	Server   config.Server   `mapstructure:"server"`
	JWT      config.JWT      `mapstructure:"jwt"`
	Crypto   config.Crypto   `mapstructure:"crypto"`
}

func New() (Config, error) {
	cfg := Config{}

	err := config.Load(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("read config error: %w", err)
	}

	return cfg, nil
}
