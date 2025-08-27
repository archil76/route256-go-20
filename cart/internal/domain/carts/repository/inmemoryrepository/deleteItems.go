package inmemoryrepository

import (
	"context"
	"route256/cart/internal/domain/model"
)

func (r *Repository) DeleteItems(ctx context.Context, userID model.UserID) (model.UserID, error) {

	cart, err := r.GetCart(ctx, userID)

	if err != nil {

		return userID, nil

	}

	clear(cart.Items)

	return userID, nil
}
