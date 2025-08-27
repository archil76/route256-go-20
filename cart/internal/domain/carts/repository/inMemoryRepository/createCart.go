package inMemoryRepository

import (
	"context"

	"route256/cart/internal/domain/model"
)

func (r *Repository) CreateCart(_ context.Context, cart model.Cart) (*model.Cart, error) {

	if cart.UserID < 1 {

		return nil, ErrUserIDIsnotValid
	}

	r.storage[cart.UserID] = cart

	return &cart, nil
}
