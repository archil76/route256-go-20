package inmemoryrepository

import (
	"errors"
	"route256/loms/internal/domain/model"
	"sync"
)

var (
	ErrOrderDoesntExist = errors.New("order doesn't exist")
	ErrUserIDIsNotValid = errors.New("UserID should be more than 0")
)

type Storage = map[int64]model.Order

type SequenceGenerator interface {
	Add(delta int64) (new int64)
}

type Repository struct {
	storage           Storage
	mu                sync.RWMutex
	sequenceGenerator SequenceGenerator
}

func NewOrderInMemoryRepository(capacity int, sequenceGenerator SequenceGenerator) *Repository {
	return &Repository{storage: make(Storage, capacity), sequenceGenerator: sequenceGenerator}
}
