package service

import (
	"context"
	"errors"
	"route256/loms/internal/domain/model"
)

var (
	ErrUserIDIsNotValid   = errors.New("user ID should be more than 0")
	ErrOrderIDIsNotValid  = errors.New("order ID should be more than 0")
	ErrSkuIDIsNotValid    = errors.New("sku should be more than 0")
	ErrOrderDoesntExist   = errors.New("order doesn't exist")
	ErrInvalidOrderStatus = errors.New("order status should be PAYED")
	ErrShortOfStock       = errors.New("available amount of stock isn't enough ")
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

type LomsService struct {
	orderRepository OrderRepository
	stockRepository StockRepository
}

func NewLomsService(orderRepository OrderRepository, stockRepository StockRepository) *LomsService {
	return &LomsService{orderRepository: orderRepository, stockRepository: stockRepository}
}
