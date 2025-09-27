package service

import (
	"context"
	"route256/cart/internal/domain/model"
)

func (s *CartService) DeleteItemByUserID(ctx context.Context, userID model.UserID) (model.UserID, error) {
	if userID < 1 {
		return 0, ErrFailValidation
	}

	_, err := s.repository.DeleteItems(ctx, userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}
