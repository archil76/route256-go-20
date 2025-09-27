package service

import (
	"context"
	"route256/cart/internal/domain/model"
)

func (s *CartService) DeleteItem(ctx context.Context, userID model.UserID, skuID model.Sku) (model.Sku, error) {
	if userID < 1 || skuID < 1 {
		return 0, ErrFailValidation
	}

	_, err := s.repository.DeleteItem(ctx, userID, model.Item{Sku: skuID, Count: 0})
	if err != nil {
		return 0, err
	}

	return skuID, nil
}
