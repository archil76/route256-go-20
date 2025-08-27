package service

import (
	"context"
	"errors"
	"route256/cart/internal/domain/model"
)

var (
	ErrInvalidSKU     = errors.New("invalid sku")
	ErrFailValidation = errors.New("fail validation")
	ErrCartIsEmpty    = errors.New("cart is empty")
)

type CartsRepository interface {
	AddItem(_ context.Context, userID model.UserID, item model.Item) (*model.Item, error)
	DeleteItem(_ context.Context, userID model.UserID, item model.Item) (*model.Item, error)
	GetCart(_ context.Context, userID model.UserID) (*model.Cart, error)
	DeleteItems(_ context.Context, userID model.UserID) (model.UserID, error)
}

type ProductService interface {
	GetProductBySku(ctx context.Context, sku model.Sku) (*model.Product, error)
}

type CartService struct {
	repository     CartsRepository
	productService ProductService
}

func NewCartsService(repository CartsRepository, productService ProductService) *CartService {
	return &CartService{repository: repository, productService: productService}
}
