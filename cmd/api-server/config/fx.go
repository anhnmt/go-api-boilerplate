package config

import (
	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/casbin"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/otel"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/redis"
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

func serverConfig(c Config) config.Server {
	return c.Server
}

func grpcConfig(c Config) config.Grpc {
	return c.Server.Grpc
}

func cryptoConfig(c Config) config.Crypto {
	return c.Crypto
}

func jwtConfig(c Config) config.JWT {
	return c.JWT
}

func postgresConfig(c Config) postgres.Config {
	return c.Postgres
}

func redisConfig(c Config) redis.Config {
	return c.Redis
}

func otelConfig(c Config) otel.Config {
	return c.Otel
}

func casbinConfig(c Config) casbin.Config {
	return c.Casbin
}

var Module = fx.Module("config", fx.Provide(
	New,
	appConfig,
	loggerConfig,
	serverConfig,
	grpcConfig,
	cryptoConfig,
	jwtConfig,
	postgresConfig,
	redisConfig,
	otelConfig,
	casbinConfig,
))
