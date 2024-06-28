package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
)

func New(ctx context.Context, cfg config.Redis) (redis.UniversalClient, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:            cfg.Address,
		Password:         cfg.Password,
		DB:               cfg.DB,
		PoolSize:         cfg.PoolSize,
		MinIdleConns:     cfg.MinIdleConns,
		MaxIdleConns:     cfg.MaxIdleConns,
		MaxActiveConns:   cfg.MaxActiveConns,
		DisableIndentity: true,
	})

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
