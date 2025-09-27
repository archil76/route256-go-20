package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) GetByID(ctx context.Context, orderID int64) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.getOrder(ctx, orderID)
}

func (r *Repository) getOrder(_ context.Context, orderID int64) (*model.Order, error) {

	if orderID < 1 {
		return nil, ErrUserIDIsNotValid
	}

	stock, ok := r.storage[orderID]
	if !ok {
		return nil, ErrOrderDoesntExist
	}
	return &stock, nil
}
