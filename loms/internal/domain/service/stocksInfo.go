package service

import (
	"context"
	"fmt"
)

func (s *LomsService) StocksInfo(ctx context.Context, sku int64) (uint32, error) {
	if sku < 1 {
		return 0, ErrSkuIDIsNotValid
	}

	var count uint32
	var err error
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		count, err = s.stockRepository.GetBySKU(ctx, sku)
		if err != nil {
			return fmt.Errorf("stockRepository.GetBySKU: %w", err)
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}
