package inmemoryrepository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewOrderPostgresRepository(pool *pgxpool.Pool) (*Repository, error) {
	repository := &Repository{
		pool: pool,
	}

	return repository, nil
}
