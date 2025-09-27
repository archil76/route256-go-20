package inmemoryrepository

import (
	"context"
	"route256/cart/internal/domain/model"
)

func (r *Repository) DeleteItems(ctx context.Context, userID model.UserID) (model.UserID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cart, err := r.getCart(ctx, userID)
	if err != nil {

		return userID, nil

	}

	clear(cart.Items)

	return userID, nil
}
