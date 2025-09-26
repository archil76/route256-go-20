package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewStockPostgresRepository(pool *pgxpool.Pool) (*Repository, error) {
	repository := &Repository{
		pool: pool,
	}

	return repository, nil
}
