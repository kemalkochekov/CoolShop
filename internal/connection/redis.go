package connection

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"Coolshop/internal/config"
)

var _ Cache = (*Redis)(nil)

type Redis struct {
	redis *redis.Client
}

func NewCache(ctx context.Context, cfgs config.RedisConfig) (*Redis, error) {
	dsn := cfgs.GenerateDSN()
	options := &redis.Options{
		Addr:     dsn,
		Password: cfgs.Password,
		DB:       0, // use default DB
	}
	rdb := redis.NewClient(options)

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("rdb.Ping: %w", err)
	}

	return &Redis{redis: rdb}, nil
}

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Close() error
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}

func (r *Redis) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.redis.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}

func (r *Redis) Close() error {
	return r.redis.Close()
}
