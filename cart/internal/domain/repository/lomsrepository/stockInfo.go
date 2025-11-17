package lomsrepository

import (
	"context"
	desc "route256/cart/internal/api"

	"route256/cart/internal/domain/model"
)

func (s *LomsService) StockInfo(ctx context.Context, sku model.Sku) (uint32, error) {

	stocksInfoRequest := desc.StocksInfoRequest{Sku: sku}

	stocksInfoResponse, err := s.client.StocksInfo(ctx, &stocksInfoRequest)
	if err != nil {
		return 0, ErrSkuNotFoundInStock
	}

	return stocksInfoResponse.Count, nil
}
