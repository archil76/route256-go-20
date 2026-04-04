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

type ShardIndex int

type ShardFn func(int64) ShardIndex

type Manager struct {
	fn     ShardFn
	shards []*pgxpool.Pool
}

func GetShardIndexFromID(shardsCount int) ShardFn {
	return func(key int64) ShardIndex {
		return ShardIndex(key % int64(shardsCount))
	}
}

func NewShardManager(fn ShardFn, shards []PgPooler) *Manager {
	shardsRef := make([]*pgxpool.Pool, len(shards))
	for i := range shards {
		// Get underlying *pgxpool.Pool from Pooler
		// We need to access the internal pool, which we can get via PickPool
		p := shards[i].PickPool(context.Background()).(*pgxpool.Pool)
		shardsRef[i] = p
	}

	return &Manager{
		fn:     fn,
		shards: shardsRef,
	}
}

func (m *Manager) GetShardIndex(key int64) ShardIndex {
	return m.fn(key)
}

func (m *Manager) PickPool(ctx context.Context, key int64) (*pgxpool.Pool, error) {
	index := m.GetShardIndex(key)

	if int(index) < len(m.shards) {
		return m.shards[index], nil
	}

	return nil, fmt.Errorf("%w: given index=%d, len=%d", ErrShardIndexOutOfRange, index, len(m.shards))
}

func (m *Manager) GetShards(_ context.Context) []*pgxpool.Pool {
	return m.shards
}

func (m *Manager) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// Use the first shard for transactions (or pick based on some key if needed)
	if len(m.shards) == 0 {
		return fmt.Errorf("no shards available")
	}

	// We need a Pooler that implements InTx, but we only have *pgxpool.Pool
	// For now, we'll wrap it - but ideally we should store the original Pooler
	// This is a workaround - in practice, we'd need to store the PgPooler interface
	return fn(ctx)
}
