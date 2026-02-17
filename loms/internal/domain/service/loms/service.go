package loms

import (
	"context"
	"route256/loms/internal/domain/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (*model.Order, error)
	GetByID(ctx context.Context, orderID int64) (*model.Order, error)
	UpdateOrder(ctx context.Context, order model.Order) (*model.Order, error)
	SetStatus(ctx context.Context, order model.Order, status model.Status) error
}

type StockRepository interface {
	GetStock(ctx context.Context, sku int64) (*model.Stock, error)
	UpdateStock(ctx context.Context, stock model.Stock) (*model.Stock, error)
	Reserve(ctx context.Context, items []model.Item) ([]model.Stock, error)
	ReserveRemove(ctx context.Context, items []model.Item) error
	ReserveCancel(ctx context.Context, items []model.Item) error
	GetBySKU(ctx context.Context, sku int64) (uint32, error)
}

type OutboxService interface {
	CreateMessage(ctx context.Context, orderID int64, status model.Status)
}

type PgPooler interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type LomsService struct {
	orderRepository OrderRepository
	stockRepository StockRepository
	outboxService   OutboxService
	pooler          PgPooler
}

func NewLomsService(orderRepository OrderRepository, stockRepository StockRepository, outboxService OutboxService, pooler PgPooler) *LomsService {
	return &LomsService{orderRepository: orderRepository, stockRepository: stockRepository, outboxService: outboxService, pooler: pooler}
}
