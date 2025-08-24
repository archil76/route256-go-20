package repository

import (
	"context"
	"errors"
	model "route256/cart/internal/domain/model"
)

var (
	ErrCartDoesntExist  = errors.New("cart doesn't exist")
	ErrUserIDIsnotValid = errors.New("UserID should be more than 0")
)

type Storage = map[model.UserID]model.Cart

type Repository struct {
	storage Storage
}

func NewCartInMemoryRepository(capacity int) *Repository {
	return &Repository{storage: make(Storage, capacity)}
}

func (r *Repository) AddItem(ctx context.Context, userID model.UserID, item model.Item) (*model.Item, error) {

	cart, err := r.GetCart(ctx, userID)

	if err != nil {

		cart, err = r.CreateCart(ctx, model.Cart{UserID: userID, Items: map[model.Sku]uint16{}})
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

func (r *Repository) CreateCart(_ context.Context, cart model.Cart) (*model.Cart, error) {

	if cart.UserID < 1 {

		return nil, ErrUserIDIsnotValid
	}

	r.storage[cart.UserID] = cart

	return &cart, nil
}

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

func (r *Repository) DeleteItem(ctx context.Context, userID model.UserID, item model.Item) (*model.Item, error) {

	cart, err := r.GetCart(ctx, userID)

	if err != nil {

		return nil, err

	}

	delete(cart.Items, item.Sku)

	return &item, nil

}

func (r *Repository) DeleteItems(ctx context.Context, userID model.UserID) (model.UserID, error) {

	cart, err := r.GetCart(ctx, userID)

	if err != nil {

		return userID, nil

	}

	clear(cart.Items)

	return userID, nil
}
