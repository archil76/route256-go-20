package postgres

import (
	manager "route256/loms/internal/infra/postgres/transactions-manager"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool      *pgxpool.Pool
	txManager *manager.TxManager
}

func NewOrderPostgresRepository(pool *pgxpool.Pool, txManager *manager.TxManager) (*Repository, error) {
	repository := &Repository{
		pool:      pool,
		txManager: txManager,
	}

	return repository, nil
}
