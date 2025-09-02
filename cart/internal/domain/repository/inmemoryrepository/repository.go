package inmemoryrepository

import (
	"errors"
	model "route256/cart/internal/domain/model"
	"sync"
)

var (
	//ErrInvalidSKU       = errors.New("invalid sku")
	ErrCartDoesntExist  = errors.New("cart doesn't exist")
	ErrUserIDIsNotValid = errors.New("UserID should be more than 0")
)

type Storage = map[model.UserID]model.Cart

type Repository struct {
	storage Storage
	mu      sync.RWMutex
}

func NewCartInMemoryRepository(capacity int) *Repository {
	return &Repository{storage: make(Storage, capacity)}
}
