package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/cockroachdb/errors"
)

func (r *Repository) Create(ctx context.Context, order model.Order) (*model.Order, error) {
	const queryOrders = `INSERT INTO orders (user_id, status) VALUES ($1, $2) returning id`
	if err := r.pool.QueryRow(ctx, queryOrders, order.UserID, order.Status).
		Scan(&order.OrderID); err != nil {
		return nil, errors.Wrap(err, "pgx.QueryRow.Scan")
	}

	const queryOrderItems = `INSERT INTO order_items (order_id, sku, count) VALUES ($1, $2, $3)`
	for _, item := range order.Items {
		if _, err := r.pool.Query(ctx, queryOrderItems, order.OrderID, item.Sku, item.Count); err != nil {
			return nil, errors.Wrap(err, "pgx.QueryRow.Scan")
		}
	}

	return &order, nil
}
