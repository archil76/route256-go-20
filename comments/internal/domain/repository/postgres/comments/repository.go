package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShardManager interface {
	PickPool(ctx context.Context, key int64) (*pgxpool.Pool, error)
	GetShards(ctx context.Context) []*pgxpool.Pool
}

type Repository struct {
	sm ShardManager
}

func NewCommentsPostgresRepository(sm ShardManager) (*Repository, error) {
	repository := &Repository{
		sm: sm,
	}

	return repository, nil
}
