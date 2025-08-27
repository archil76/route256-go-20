package inMemoryRepository

import (
	"context"

	"route256/cart/internal/domain/model"
)

func (r *Repository) DeleteItem(ctx context.Context, userID model.UserID, item model.Item) (*model.Item, error) {

	cart, err := r.GetCart(ctx, userID)

	if err != nil {

		return nil, err

	}

	delete(cart.Items, item.Sku)

	return &item, nil

}
