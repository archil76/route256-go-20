package inmemoryrepository

import (
	"context"
	"route256/cart/internal/domain/model"
)

func (r *Repository) GetCart(ctx context.Context, userID model.UserID) (*model.Cart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.getCart(ctx, userID)
}

func (r *Repository) getCart(_ context.Context, userID model.UserID) (*model.Cart, error) {
	if userID < 1 {
		return nil, ErrUserIDIsNotValid
	}

	cart, ok := r.storage[userID]
	if !ok {
		return nil, ErrCartDoesntExist
	}
	return &cart, nil
}
