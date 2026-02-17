package inmemoryrepository

import (
	"context"
	"errors"
	"route256/cart/internal/domain/model"
	"route256/cart/internal/infra/logger"
	"route256/cart/internal/infra/metrics"
	"sync"
	"time"

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
	done    chan struct{}
	tracer  Tracer
}

func NewCartInMemoryRepository(capacity int, tracer Tracer) *Repository {
	repository := &Repository{storage: make(Storage, capacity), done: make(chan struct{}), tracer: tracer}
	go func() {
		t := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-t.C:
				repository.mu.Lock()
				storageLen := len(repository.storage)
				metrics.StoreRepoSize(float64(storageLen))
				logger.Infow("Repo size: ", "storageLen", storageLen)
				repository.mu.Unlock()
			case <-repository.done:
				t.Stop()
				return
			}
		}
	}()
	return repository
}
