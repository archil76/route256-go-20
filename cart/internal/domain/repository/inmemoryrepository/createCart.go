package inmemoryrepository

import (
	"context"

	"route256/cart/internal/domain/model"
)

func (r *Repository) createCart(_ context.Context, cart model.Cart) (*model.Cart, error) {
	if cart.UserID < 1 {

		return nil, ErrUserIDIsNotValid
	}

	r.storage[cart.UserID] = cart

	return &cart, nil
}
