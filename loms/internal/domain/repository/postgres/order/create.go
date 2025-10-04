package postgres

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Create(ctx context.Context, order model.Order) (*model.Order, error) {
	var err error

	txManager := r.txManager
	pool := txManager.GetPool()

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return nil, err
	}

	//ctx = ctxWithTx(ctx, tx)

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

	//err = fn(ctx)

	const queryOrders = `INSERT INTO orders (user_id, status) VALUES ($1, $2) returning id`
	if err = tx.QueryRow(ctx, queryOrders, order.UserID, order.Status).
		Scan(&order.OrderID); err != nil {
		err = errors.Wrap(err, "pgx.QueryRow.Scan")
		return nil, err
	}

	const queryOrderItems = `INSERT INTO order_items (order_id, sku, count) VALUES ($1, $2, $3)`
	batch := &pgx.Batch{}
	for _, item := range order.Items {
		batch.Queue(queryOrderItems, order.OrderID, item.Sku, item.Count)
	}
	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err = br.Exec()
		if err != nil {
			err = errors.Wrap(err, "pgx.QueryRow.Scan")
			return nil, err
		}
	}

	return &order, nil
}
