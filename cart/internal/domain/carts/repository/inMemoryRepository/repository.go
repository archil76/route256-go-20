package inMemoryRepository

import (
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
