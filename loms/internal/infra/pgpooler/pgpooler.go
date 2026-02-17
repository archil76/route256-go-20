package pgpooler

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Pooler struct {
		master *pgxpool.Pool
	}

	Pool interface {
		Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
		QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
		SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	}

	ctxTxType struct{}
)

func NewPooler(ctx context.Context, dsn string) (*Pooler, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.ParseConfig")
	}

	config.ConnConfig.Tracer = NewPgxQueryTracer()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.NewWithConfig")
	}

	pooler := &Pooler{master: pool}

	return pooler, nil
}

func (p *Pooler) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := p.master.Begin(ctx)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, ctxTxType{}, tx)

	err = fn(ctx)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return fmt.Errorf("failed to rollback transaction after error: %w, rollback error: %w", err, txErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

func (p *Pooler) PickPool(ctx context.Context) Pool {
	if tx := ctx.Value(ctxTxType{}); tx != nil {
		return tx.(Pool)
	}
	return p.master
}
