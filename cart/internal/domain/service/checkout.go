package service

import (
	"context"
	"route256/cart/internal/domain/model"
)

func (s *CartService) Checkout(ctx context.Context, userID model.UserID) (int64, error) {
	if userID < 1 {
		return 0, ErrFailValidation
	}

	reportCart, err := s.GetItemsByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	orderID, err := s.lomsService.OrderCreate(ctx, userID, reportCart)
	if err != nil && orderID != 0 {
		return 0, err
	}

	_, err = s.DeleteItemByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}
