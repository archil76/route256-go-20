package inmemoryrepository

import (
	"context"
	"route256/cart/internal/domain/model"
)

func (r *Repository) AddItem(ctx context.Context, userID model.UserID, item model.Item) (*model.Item, error) {

	cart, err := r.GetCart(ctx, userID)

	if err != nil {

		cart, err = r.CreateCart(ctx, model.Cart{UserID: userID, Items: map[model.Sku]uint32{}})
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
