package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jesusthecreator017/PicoURL/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// Create a new redis cache instance
func NewRedisCache(cfg *config.RedisConfig, ttl time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &RedisCache{client: client, ttl: ttl}, nil
}

// Basic Redis operations for caching URL mappings
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key, value string) error {
	return r.client.Set(ctx, key, value, r.ttl).Err()
}

func (r *RedisCache) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisCache) Incr(ctx context.Context, key string) error {
	return r.client.Incr(ctx, key).Err()
}

func (r *RedisCache) GetInt(ctx context.Context, key string) (int, error) {
	count, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, redis.Nil // Signal a Cache Miss
	}
	return int(count), err
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
