package service

import (
	"context"
	"errors"
	"route256/cart/internal/domain/model"
)

var (
	ErrFailValidation = errors.New("fail validation")
	ErrCartIsEmpty    = errors.New("cart is empty")
)

type CartsRepository interface {
	AddItem(ctx context.Context, userID model.UserID, item model.Item) (*model.Item, error)
	DeleteItem(ctx context.Context, userID model.UserID, item model.Item) (*model.Item, error)
	GetCart(ctx context.Context, userID model.UserID) (*model.Cart, error)
	DeleteItems(ctx context.Context, userID model.UserID) (model.UserID, error)
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
