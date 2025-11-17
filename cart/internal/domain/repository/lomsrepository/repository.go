package lomsrepository

import (
	"errors"

	desc "route256/cart/internal/api"

	"route256/cart/internal/infra/http/middlewares"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrSkuNotFoundInStock = errors.New("sku not found in stock")
	ErrOrderNotFound      = errors.New("product not found")
)

type LomsService struct {
	address string
	client  desc.LomsClient
}

func NewLomsService(address string) (*LomsService, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middlewares.CounterUnaryClientInterceptor),
	)
	if err != nil {
		return nil, err
	}

	client := desc.NewLomsClient(conn)

	return &LomsService{
		address: address,
		client:  client,
	}, nil
}
