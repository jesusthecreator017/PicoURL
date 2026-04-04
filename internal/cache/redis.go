package store

import (
	"context"
	"fmt"

	"github.com/jesusthecreator017/PicoURL/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

// Create a new redis store instance
func NewRedisStore(cfg *config.RedisConfig) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &RedisStore{client: client}, nil
}

// Store Methods
func (r *RedisStore) SaveURL(shortURL, originalURL string) error {
	return r.client.Set(context.Background(), shortURL, originalURL, 0).Err()
}

func (r *RedisStore) GetOriginalURL(shortURL string) (string, error) {
	return r.client.Get(context.Background(), shortURL).Result()
}

func (r *RedisStore) IncrementCount(shortURL string) error {
	return r.client.Incr(context.Background(), shortURL+":count").Err()
}

func (r *RedisStore) GetCount(shortURL string) (int, error) {
	count, err := r.client.Get(context.Background(), shortURL+":count").Int()
	if err == redis.Nil {
		return 0, nil // No count exists yet
	}
	return count, err
}
