package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) Create(_ context.Context, order model.Order) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order.OrderId = r.sequenceGenerator.Add(1)

	r.storage[order.OrderId] = order

	return &order, nil
}
