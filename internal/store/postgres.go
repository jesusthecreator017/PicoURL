package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jesusthecreator017/PicoURL/internal/config"
	"github.com/jesusthecreator017/PicoURL/internal/store/sqlc"
)

// Create the PostgresStore struct that implements the Store interface
type PostgresStore struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// Create a new postgres store instance
func NewPostgresStore(cfg *config.PostgresConfig) (*PostgresStore, error) {
	// Build the DSN (Data Source Name) for PostgreSQL connection
	dsn := config.DSN(cfg)

	// Create a connection pool
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}

	return &PostgresStore{
		pool:    pool,
		queries: sqlc.New(pool),
	}, nil
}

func (p *PostgresStore) SaveURL(ctx context.Context, shortURL, originalURL string) error {
	return p.queries.SaveURL(ctx, sqlc.SaveURLParams{
		ShortUrl:    shortURL,
		OriginalUrl: originalURL,
	})
}

func (p *PostgresStore) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	return p.queries.GetOriginalURL(ctx, shortURL)
}

func (p *PostgresStore) IncrementCount(ctx context.Context, shortURL string) error {
	return p.queries.IncrementCount(ctx, shortURL)
}

func (p *PostgresStore) GetCount(ctx context.Context, shortURL string) (int, error) {
	count, err := p.queries.GetCount(ctx, shortURL)
	return int(count), err
}

func (p *PostgresStore) DeleteURL(ctx context.Context, shortURL string) error {
	return p.queries.DeleteURL(ctx, shortURL)
}

func (p *PostgresStore) Close() error {
	p.pool.Close()
	return nil
}
