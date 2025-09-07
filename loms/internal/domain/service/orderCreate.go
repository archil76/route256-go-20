package service

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (s *LomsService) OrderCreate(ctx context.Context, userID int64, items []model.Item) (int64, error) {
	if userID < 1 {
		return 0, ErrUserIDIsNotValid
	}

	order := model.Order{
		OrderId: 0,
		UserID:  userID,
		Status:  model.NEW_STATUS,
		Items:   items,
	}

	upOrder, err := s.orderRepository.Create(ctx, order)
	if err != nil {
		return 0, err
	}

	orderStatus := model.AWAITING_PAYMENT

	_, err = s.stockRepository.Reserve(ctx, items)
	if err != nil {
		orderStatus = model.FAILED
	}

	err = s.orderRepository.SetStatus(ctx, *upOrder, orderStatus)
	if err != nil {
		return upOrder.OrderId, err // Заказ уже записан в статусе new. Так что id можно вернуть.
	}

	return upOrder.OrderId, nil
}
