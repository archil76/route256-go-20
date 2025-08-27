package service

import (
	"context"
	"errors"
	"route256/cart/internal/domain/model"
)

func (s *CartService) AddItem(ctx context.Context, userID model.UserID, skuID model.Sku, count uint16) (model.Sku, error) {

	if userID < 1 || skuID < 1 || count < 1 {
		return 0, ErrFailValidation
	}

	if _, err := s.productService.GetProductBySku(ctx, skuID); err != nil {

		if errors.Is(err, model.ErrProductNotFound) {
			return 0, err
		}
		return 0, ErrInvalidSKU
	}

	_, err := s.repository.AddItem(ctx, userID, model.Item{Sku: skuID, Count: count})
	if err != nil {
		return 0, err
	}

	return skuID, nil
}
