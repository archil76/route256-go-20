package inmemoryrepository

import (
	"context"
	"route256/cart/internal/domain/model"
)

func (r *Repository) AddItem(ctx context.Context, userID model.UserID, item model.Item) (*model.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cart, err := r.getCart(ctx, userID)

	if err != nil {
		cart, err = r.createCart(ctx, model.Cart{UserID: userID, Items: map[model.Sku]uint32{}})
		if err != nil {
			return nil, err
		}
	}

	if _, ok := cart.Items[item.Sku]; !ok {
		cart.Items[item.Sku] = item.Count
	} else {
		cart.Items[item.Sku] += item.Count
	}

	return &item, nil
}
