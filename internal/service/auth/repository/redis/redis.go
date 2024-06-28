package authredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	TokenBlacklistKey   = "token_blacklist"
	SessionBlacklistKey = "session_blacklist"
)

type Redis struct {
	rdb redis.UniversalClient
}

func New(rdb redis.UniversalClient) *Redis {
	return &Redis{
		rdb: rdb,
	}
}

func (r *Redis) SetTokenBlacklist(ctx context.Context, tokenId, sessionId string) error {
	key := fmt.Sprintf("%s:%s", TokenBlacklistKey, tokenId)

	return r.rdb.Set(ctx, key, sessionId, 24*time.Hour).Err()
}

func (r *Redis) CheckTokenBlacklist(ctx context.Context, tokenId string) error {
	key := fmt.Sprintf("%s:%s", TokenBlacklistKey, tokenId)

	if r.rdb.Exists(ctx, key).Val() >= 1 {
		return fmt.Errorf("token was revoked")
	}

	return nil
}

func (r *Redis) SetSessionBlacklist(ctx context.Context, sessionId string) error {
	key := fmt.Sprintf("%s:%s", SessionBlacklistKey, sessionId)

	return r.rdb.Set(ctx, key, sessionId, 24*time.Hour).Err()
}

func (r *Redis) CheckSessionBlacklist(ctx context.Context, sessionId string) error {
	key := fmt.Sprintf("%s:%s", SessionBlacklistKey, sessionId)

	if r.rdb.Exists(ctx, key).Val() >= 1 {
		return fmt.Errorf("session was revoked")
	}

	return nil
}
