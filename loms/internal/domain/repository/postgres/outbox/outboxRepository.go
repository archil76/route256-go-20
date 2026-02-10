package postgres

import (
	"context"
	"route256/loms/internal/infra/pgpooler"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Pool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}
type PgPooler interface {
	PickPool(ctx context.Context) Pool
}

type Repository struct {
	pooler *pgpooler.Pooler
}

func NewOutboxPostgresRepository(pooler *pgpooler.Pooler) (*Repository, error) {
	repository := &Repository{
		pooler: pooler,
	}

	return repository, nil
}
