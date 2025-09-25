package inmemoryrepository

import (
	"route256/loms/internal/domain/model"
	"sync"
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
