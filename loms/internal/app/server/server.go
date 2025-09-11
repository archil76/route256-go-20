package server

import (
	"context"
	desc "route256/loms/internal/api"
	"route256/loms/internal/domain/model"
)

type LomsServise interface {
	OrderCancel(ctx context.Context, orderID int64) error
	OrderCreate(ctx context.Context, userID int64, items []model.Item) (int64, error)
	OrderInfo(ctx context.Context, orderID int64) (*model.Order, error)
	OrderPay(ctx context.Context, orderId int64) error
	StocksInfo(ctx context.Context, sku int64) (uint32, error)
}

type Server struct {
	lomsServise LomsServise
	desc.UnimplementedLomsServer
}

func NewServer(lomsServise LomsServise) *Server {
	return &Server{lomsServise: lomsServise}
}
