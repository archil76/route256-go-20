package shard_manager

import (
	"context"
	"errors"
	"fmt"
	"route256/comments/internal/infra/pgpooler"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrShardIndexOutOfRange = errors.New("shard index is out of range")
)

type Pool = pgpooler.Pool

type PgPooler interface {
	PickPool(ctx context.Context) Pool
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type ShardFn func(int64) int

type Manager struct {
	fn     ShardFn
	shards []*pgxpool.Pool
}

func GetShardIndexFromID(shardsCount int) ShardFn {
	return func(key int64) int {
		return int(key % int64(shardsCount))
	}
}

func NewShardManager(fn ShardFn, shards []PgPooler) *Manager {
	shardsRef := make([]*pgxpool.Pool, len(shards))
	for i := range shards {
		p := shards[i].PickPool(context.Background()).(*pgxpool.Pool)
		shardsRef[i] = p
	}

	return &Manager{
		fn:     fn,
		shards: shardsRef,
	}
}

func (m *Manager) GetShardIndex(key int64) int {
	return m.fn(key)
}

func (m *Manager) GetShardIndexFromID(id int64) int {
	// 123002
	// 123- seq с шарда
	// 2- номер шарда
	return int(id % 1000)
}

func (m *Manager) PickPool(_ context.Context, shardIndex int) (*pgxpool.Pool, error) {

	if int(shardIndex) < len(m.shards) {
		return m.shards[shardIndex], nil
	}

	return nil, fmt.Errorf("%w: given index=%d, len=%d", ErrShardIndexOutOfRange, shardIndex, len(m.shards))
}

func (m *Manager) GetShards(_ context.Context) []*pgxpool.Pool {
	return m.shards
}

func (m *Manager) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if len(m.shards) == 0 {
		return fmt.Errorf("no shards available")
	}

	return fn(ctx)
}
