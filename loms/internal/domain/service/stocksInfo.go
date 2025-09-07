package service

import "context"

func (s *LomsService) StocksInfo(ctx context.Context, sku int64) (uint32, error) {
	if sku < 1 {
		return 0, ErrSkuIDIsNotValid
	}

	count, err := s.stockRepository.GetBySKU(ctx, sku)
	if err != nil {
		return 0, err
	}

	return count, nil
}
