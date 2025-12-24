package database

import (
	"context"
	"fmt"

	"route256/tests/app/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager interface {
	RunMaster(ctx context.Context, fn func(ctxTx context.Context, tx Transaction) error) error
	RunReplica(ctx context.Context, fn func(ctxTx context.Context, tx Transaction) error) error
	Conn() Transaction
}

type Transaction interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type PgTxManager struct {
	poolMaster  *pgxpool.Pool
	poolReplica *pgxpool.Pool
}

func NewPgTxManager(poolMaster, poolReplica *pgxpool.Pool) *PgTxManager {
	return &PgTxManager{
		poolMaster:  poolMaster,
		poolReplica: poolReplica,
	}
}

func (m *PgTxManager) RunMaster(ctx context.Context, fn func(ctxTx context.Context, tx Transaction) error) error {
	options := pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	}
	return m.inTx(ctx, m.poolMaster, options, fn)
}

func (m *PgTxManager) RunReplica(ctx context.Context, fn func(ctxTx context.Context, tx Transaction) error) error {
	options := pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	}
	return m.inTx(ctx, m.poolReplica, options, fn)
}

func (m *PgTxManager) IsReplicaEnabled() bool {
	return m.poolReplica != nil
}

func (m *PgTxManager) Conn() Transaction {
	return m.poolMaster
}

func (m *PgTxManager) inTx(ctx context.Context, pool *pgxpool.Pool, options pgx.TxOptions, f func(ctxTx context.Context, tx Transaction) error) error {
	tx, err := pool.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to begin tx, err: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			logger.Info("%v", p)
			_ = tx.Rollback(ctx)
			panic(p) // fallthrough panic after rollback on caught panic
		} else if err != nil {
			_ = tx.Rollback(ctx) // if error during computations
		} else {
			err = tx.Commit(ctx) // all good
		}
	}()

	return f(ctx, tx)
}
