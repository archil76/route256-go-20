package postgres

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Create(ctx context.Context, order model.Order) (*model.Order, error) {
	var err error

	err = r.pooler.InTx(ctx, "R", func(ctx context.Context) error {
		pool := r.pooler.PickPool(ctx)

		const queryOrders = `INSERT INTO orders (user_id, status) VALUES ($1, $2) returning id`
		if err = pool.QueryRow(ctx, queryOrders, order.UserID, order.Status).
			Scan(&order.OrderID); err != nil {
			err = errors.Wrap(err, "pgx.QueryRow.Scan")
			return err
		}

		const queryOrderItems = `INSERT INTO order_items (order_id, sku, count) VALUES ($1, $2, $3)`
		batch := &pgx.Batch{}
		for _, item := range order.Items {
			batch.Queue(queryOrderItems, order.OrderID, item.Sku, item.Count)
		}
		br := pool.SendBatch(ctx, batch)
		defer br.Close()

		//for i := 0; i < batch.Len(); i++ {
		//	_, err = br.Exec()
		//	if err != nil {
		//		err = errors.Wrap(err, "pgx.QueryRow.Scan")
		//		return err
		//	}
		//}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &order, nil
}
