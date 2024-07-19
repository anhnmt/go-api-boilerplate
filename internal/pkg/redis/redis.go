package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Ctx    context.Context
	Config Config
}

func New(lc fx.Lifecycle, p Params) (redis.UniversalClient, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:            p.Config.Address,
		Password:         p.Config.Password,
		DB:               p.Config.DB,
		PoolSize:         p.Config.PoolSize,
		MinIdleConns:     p.Config.MinIdleConns,
		MaxIdleConns:     p.Config.MaxIdleConns,
		MaxActiveConns:   p.Config.MaxActiveConns,
		DisableIndentity: true,
	})

	ctx, cancel := context.WithTimeout(p.Ctx, 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	lc.Append(fx.StopHook(func() error {
		return client.Close()
	}))

	return client, nil
}
