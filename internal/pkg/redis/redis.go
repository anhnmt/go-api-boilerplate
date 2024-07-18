package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func New(ctx context.Context, lc fx.Lifecycle, cfg Config) (redis.UniversalClient, error) {
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

	lc.Append(fx.StopHook(func() error {
		return client.Close()
	}))

	return client, nil
}
