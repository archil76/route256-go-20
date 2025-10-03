package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/cockroachdb/errors"
)

func (r *Repository) GetByID(ctx context.Context, orderID int64) (*model.Order, error) {
	return r.getOrder(ctx, orderID)
}

func (r *Repository) getOrder(ctx context.Context, orderID int64) (*model.Order, error) {

	if orderID < 1 {
		return nil, model.ErrUserIDIsNotValid
	}

	const queryOrders = `SELECT id, user_id, status FROM orders where id = $1`

	upOrder := model.Order{}

	if err := r.pool.QueryRow(ctx, queryOrders, orderID).
		Scan(&upOrder.OrderID, &upOrder.UserID, &upOrder.Status); err != nil {
		return nil, errors.Wrap(model.ErrOrderDoesntExist, "pgx.QueryRow.Scan")
	}

	const queryOrderItems = `SELECT sku, count FROM order_items where order_id = $1`

	rows, err := r.pool.Query(ctx, queryOrderItems, orderID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		upItem := model.Item{}
		if err := rows.Scan(&upItem.Sku, &upItem.Count); err != nil {
			return nil, err
		}
		upOrder.Items = append(upOrder.Items, model.Item{
			Sku:   upItem.Sku,
			Count: upItem.Count,
		})
	}

	return &upOrder, nil
}
