package inMemoryRepository

import (
	"context"
	"route256/cart/internal/domain/model"
)

func (r *Repository) GetCart(_ context.Context, userID model.UserID) (*model.Cart, error) {

	if userID < 1 {
		return nil, ErrUserIDIsnotValid
	}

	cart, ok := r.storage[userID]
	if !ok {
		return nil, ErrCartDoesntExist
	}
	return &cart, nil
}
