package server

import (
	"context"
	"errors"
	"route256/cart/internal/domain/model"
)

var (
	ErrInvalidUserID = errors.New("Идентификатор пользователя должен быть натуральным числом (больше нуля)")
	ErrInvalidSKU    = errors.New("SKU должен быть натуральным числом (больше нуля)")
	ErrInvalidCount  = errors.New("количество должно быть натуральным числом (больше нуля)")
	ErrPSFail        = errors.New("SKU должен существовать в сервисе product-service")
	ErrUnmarshalling = errors.New("Unmarshalling error")
	ErrOther         = errors.New("Ошибка сервера")
)

type CartService interface {
	AddItem(ctx context.Context, userID model.UserID, skuID model.Sku, count uint32) (model.Sku, error)
	DeleteItem(ctx context.Context, userID model.UserID, skuID model.Sku) (model.Sku, error)
	DeleteItemByUserID(ctx context.Context, userID model.UserID) (model.UserID, error)
	GetItemsByUserID(ctx context.Context, userID model.UserID) (*model.ReportCart, error)
}

type Server struct {
	cartService CartService
}

func NewServer(cartService CartService) *Server {
	return &Server{cartService: cartService}
}
