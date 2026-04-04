package store

import "github.com/jackc/pgx/v5/pgxpool"

type PostgresStore struct {
	pool *pgxpool.Pool
}
