package cache

import (
	"Coolshop/internal/connection"
	"Coolshop/pkg/constant"
	"context"
	"fmt"
	"time"
)

type Repository interface {
	SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
}

var _ Repository = (*RedisRepo)(nil)

type RedisRepo struct {
	cache connection.Cache
}

// NewRedisRepo initializes RedisRepo.
func NewRedisRepo(cache connection.Cache) *RedisRepo {
	return &RedisRepo{cache: cache}
}

// SaveRefreshToken stores a refresh token in Redis with expiration.
func (r *RedisRepo) SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	expiration := time.Hour * constant.CacheExpirationTime // 7 days expiration

	err := r.cache.Set(ctx, userID, refreshToken, expiration)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

// GetRefreshToken retrieves the refresh token from Redis for a given userID.
func (r *RedisRepo) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	token, err := r.cache.Get(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}

	return token, nil
}

// DeleteRefreshToken removes the refresh token from Redis.
func (r *RedisRepo) DeleteRefreshToken(ctx context.Context, userID string) error {
	err := r.cache.Del(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}
