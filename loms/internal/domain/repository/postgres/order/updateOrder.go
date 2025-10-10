package postgres

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/cockroachdb/errors"
)

func (r *Repository) UpdateOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	return r.updateOrder(ctx, order)
}

func (r *Repository) updateOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	pool := r.pooler.PickPool(ctx)

	const queryOrders = `UPDATE orders SET user_id=$2, status=$3 where id = $1`

	upOrder := model.Order{
		OrderID: order.OrderID,
		UserID:  order.UserID,
		Status:  order.Status,
	}

	if _, err := pool.Exec(ctx, queryOrders, order.OrderID, order.UserID, order.Status); err != nil {
		return nil, errors.Wrap(err, "pgx.QueryRow.Scan")
	}

	const queryOrderItems = `UPDATE order_items SET order_id=$1, sku=$2, count=$3 where order_items.order_id= $1`
	for _, item := range order.Items {
		if _, err := pool.Exec(ctx, queryOrderItems, order.OrderID, item.Sku, item.Count); err != nil {
			return nil, errors.Wrap(err, "pgx.QueryRow.Scan")
		}
		upOrder.Items = append(upOrder.Items, model.Item{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	return &upOrder, nil
}
