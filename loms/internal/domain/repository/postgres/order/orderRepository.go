package postgres

import (
	"route256/loms/internal/infra/pgpooler"
)

type Repository struct {
	pooler *pgpooler.Pooler
}

func NewOrderPostgresRepository(pooler *pgpooler.Pooler) (*Repository, error) {
	repository := &Repository{
		pooler: pooler,
	}

	return repository, nil
}
