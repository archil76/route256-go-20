package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) UpdateOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.updateOrder(ctx, order)
}

func (r *Repository) updateOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	_, err := r.getOrder(ctx, order.OrderId)
	if err != nil {
		return nil, err
	}

	r.storage[order.OrderId] = order

	return &order, nil
}
