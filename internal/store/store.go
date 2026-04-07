package store

import "context"

type Store interface {
	SaveURL(ctx context.Context, shortURL, originalURL string) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	IncrementCount(ctx context.Context, shortURL string) error
	GetCount(ctx context.Context, shortURL string) (int, error)
	DeleteURL(ctx context.Context, shortURL string) error
	Close() error
}
