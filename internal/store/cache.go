package store

import (
	"context"
	"fmt"
)

type CachedStore struct {
	db    Store
	cache *RedisCache
}

func NewCachedStore(db Store, cache *RedisCache) *CachedStore {
	return &CachedStore{
		db:    db,
		cache: cache,
	}
}

func (s *CachedStore) SaveURL(ctx context.Context, shortURL, originalURL string) error {
	if err := s.db.SaveURL(ctx, shortURL, originalURL); err != nil {
		return err
	}

	_ = s.cache.Set(ctx, shortURL, originalURL)
	return nil
}

func (s *CachedStore) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	// Try to get from cache first
	originalURL, err := s.cache.Get(ctx, shortURL)
	if err == nil {
		return originalURL, nil // cache hit
	}
	// err is redis.Nil (cache miss) or a Redis error — either way fall through to Postgres

	// Cache miss — get from DB
	originalURL, err = s.db.GetOriginalURL(ctx, shortURL)
	if err != nil {
		return "", err
	}

	// Cache the result for future requests
	_ = s.cache.Set(ctx, shortURL, originalURL)
	return originalURL, nil
}

func (s *CachedStore) IncrementCount(ctx context.Context, shortURL string) error {
	if err := s.db.IncrementCount(ctx, shortURL); err != nil {
		return err
	}

	_ = s.cache.Incr(ctx, shortURL+":count")
	return nil
}

func (s *CachedStore) GetCount(ctx context.Context, shortURL string) (int, error) {
	// Try the cache first
	count, err := s.cache.GetInt(ctx, shortURL+":count")
	if err == nil {
		return count, nil
	}
	// err is redis.Nil (cache miss) or a Redis error — either way fall through to Postgres

	// Cache miss, get from DB
	count, err = s.db.GetCount(ctx, shortURL)
	if err != nil {
		return 0, err
	}

	// Populate cache for future requests
	_ = s.cache.Set(ctx, shortURL+":count", fmt.Sprintf("%d", count))
	return count, nil
}

func (s *CachedStore) DeleteURL(ctx context.Context, shortURL string) error {
	if err := s.db.DeleteURL(ctx, shortURL); err != nil {
		return err
	}

	_ = s.cache.Del(ctx, shortURL, shortURL+":count")
	return nil
}

func (s *CachedStore) Close() error {
	_ = s.cache.Close()
	return s.db.Close()
}
