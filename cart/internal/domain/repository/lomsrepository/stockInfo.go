package lomsrepository

import (
	"context"
	desc "route256/cart/internal/api"

	"route256/cart/internal/domain/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *LomsService) StockInfo(ctx context.Context, sku model.Sku) (uint32, error) {
	conn, err := grpc.NewClient(
		s.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return 0, err
	}

	client := desc.NewLomsClient(conn)

	stocksInfoRequest := desc.StocksInfoRequest{Sku: sku}

	stocksInfoResponse, err := client.StocksInfo(ctx, &stocksInfoRequest)
	if err != nil {
		return 0, ErrSkuNotFoundInStock
	}

	return stocksInfoResponse.Count, nil
}
