//go:build tx

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PoolManager struct{}

func (m *PoolManager) GetPool(ctx context.Context, role string) *pgxpool.Pool {
	panic("todo")
	return nil
}

// Tx транзакция.
type Tx pgx.Tx

type txKey struct{}

func ctxWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func TxFromCtx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)

	return tx, ok
}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewPgTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{
		pool: pool,
	}
}

func (m *TxManager) Repository(ctx context.Context) *repository.Repository {
	if tx, ok := TxFromCtx(ctx); ok {
		return repository.New(tx)
	}

	return repository.New(m.pool)
}

// WithTransaction выполняет fn в транзакции с дефолтным уровнем изоляции.
func (m *TxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	return m.WithTx(ctx, pgx.TxOptions{}, fn)
}

// WithTransaction выполняет fn в транзакции с уровнем изоляции RepeatableRead.
func (m *TxManager) WithRepeatableRead(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	return m.WithTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}, fn)
}

// WithTx выполняет fn в транзакции.
func (m *TxManager) WithTx(ctx context.Context, options pgx.TxOptions, fn func(ctx context.Context) error) (err error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "Transaction")
	defer span.Finish()

	tx, err := m.pool.BeginTx(ctx, options)
	if err != nil {
		return
	}
	ctx = ctxWithTx(ctx, tx)

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback(ctx)
		} else {
			// all good, commit
			err = tx.Commit(ctx)
		}
	}()

	err = fn(ctx)

	return
}
