package loms

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (s *LomsService) StocksInfo(ctx context.Context, sku int64) (uint32, error) {
	if sku < 1 {
		return 0, model.ErrSkuIDIsNotValid
	}

	count, err := s.stockRepository.GetBySKU(ctx, sku)
	if err != nil {
		return 0, err
	}

	return count, nil
}
