package inmemoryrepository

import (
	"context"
	"errors"
	"route256/cart/internal/domain/model"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

var (
	ErrCartDoesntExist  = errors.New("cart doesn't exist")
	ErrUserIDIsNotValid = errors.New("UserID should be more than 0")
)

type Storage = map[model.UserID]model.Cart

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type Repository struct {
	storage Storage
	mu      sync.RWMutex
	tracer  Tracer
}

func NewCartInMemoryRepository(capacity int, tracer Tracer) *Repository {
	return &Repository{storage: make(Storage, capacity), tracer: tracer}
}
