package server

import (
	"context"
	"route256/cart/internal/domain/model"
)

type CartService interface {
	AddItem(ctx context.Context, userID model.UserID, skuID model.Sku, count uint16) (model.Sku, error)
	DeleteItem(ctx context.Context, userID model.UserID, skuID model.Sku) (model.Sku, error)
	DeleteItemByUserId(ctx context.Context, userID model.UserID) (model.UserID, error)
	GetItemsByUserId(ctx context.Context, userID model.UserID) (*model.ReportCart, error)
}

type Server struct {
	cartService CartService
}

func NewServer(cartService CartService) *Server {
	return &Server{cartService: cartService}
}
