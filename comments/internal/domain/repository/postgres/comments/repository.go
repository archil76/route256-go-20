package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShardManager interface {
	GetShardIndexFromID(id int64) int
	GetShardIndex(key int64) int
	PickPool(_ context.Context, shardIndex int) (*pgxpool.Pool, error)
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
