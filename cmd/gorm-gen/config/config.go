package config

import (
	"fmt"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
)

type Config struct {
	Log      config.Log      `mapstructure:"log"`
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
