package authredis

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    "go.uber.org/fx"
)

const (
    TokenBlacklistKey   = "token_blacklist"
    SessionBlacklistKey = "session_blacklist"
)

type Redis struct {
    rdb redis.UniversalClient
}

type Params struct {
    fx.In

    RDB redis.UniversalClient
}

func New(p Params) *Redis {
    return &Redis{
        rdb: p.RDB,
    }
}

func (r *Redis) SetTokenBlacklist(ctx context.Context, tokenID, sessionID string) error {
    key := fmt.Sprintf("%s:%s", TokenBlacklistKey, tokenID)

    return r.rdb.Set(ctx, key, sessionID, 24*time.Hour).Err()
}

func (r *Redis) CheckTokenBlacklist(ctx context.Context, tokenID string) error {
    key := fmt.Sprintf("%s:%s", TokenBlacklistKey, tokenID)

    if r.rdb.Exists(ctx, key).Val() >= 1 {
        return fmt.Errorf("token was revoked")
    }

    return nil
}

func (r *Redis) SetSessionBlacklist(ctx context.Context, sessionID string) error {
    key := fmt.Sprintf("%s:%s", SessionBlacklistKey, sessionID)

    return r.rdb.Set(ctx, key, sessionID, 24*time.Hour).Err()
}

func (r *Redis) SetSessionsBlacklist(ctx context.Context, sessionIDs []string) error {
    _, err := r.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
        for _, sessionID := range sessionIDs {
            key := fmt.Sprintf("%s:%s", SessionBlacklistKey, sessionID)
            pipe.Set(ctx, key, sessionID, 24*time.Hour)
        }
        return nil
    })
    if err != nil {
        return err
    }

    return nil
}

func (r *Redis) CheckSessionBlacklist(ctx context.Context, sessionID string) error {
    key := fmt.Sprintf("%s:%s", SessionBlacklistKey, sessionID)

    if r.rdb.Exists(ctx, key).Val() >= 1 {
        return fmt.Errorf("session was revoked")
    }

    return nil
}
