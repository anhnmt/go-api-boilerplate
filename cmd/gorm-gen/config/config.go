package config

import (
	"fmt"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
)

type Config struct {
	Log      logger.Config   `mapstructure:"log"`
	Postgres config.Postgres `mapstructure:"postgres"`
}

func New() (Config, error) {
	cfg := Config{}

	err := config.Load(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("read config error: %w", err)
	}

	return cfg, nil
}
