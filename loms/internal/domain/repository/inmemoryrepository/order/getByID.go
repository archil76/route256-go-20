package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) GetByID(ctx context.Context, orderId int64) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.getOrder(ctx, orderId)
}

func (r *Repository) getOrder(_ context.Context, orderId int64) (*model.Order, error) {

	if orderId < 1 {
		return nil, ErrUserIDIsNotValid
	}

	stock, ok := r.storage[orderId]
	if !ok {
		return nil, ErrOrderDoesntExist
	}
	return &stock, nil
}
