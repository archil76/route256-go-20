package loms

import (
	"context"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/infra/logger"
)

func (s *LomsService) OrderCreate(ctx context.Context, userID int64, items []model.Item) (int64, error) {
	if userID < 1 {
		return 0, model.ErrUserIDIsNotValid
	}

	order := model.Order{
		OrderID: 0,
		UserID:  userID,
		Status:  model.NEWSTATUS,
		Items:   items,
	}

	upOrder, err := s.orderRepository.Create(ctx, order)
	if err != nil {
		return 0, err
	}

	s.outboxService.CreateMessage(ctx, upOrder.OrderID, upOrder.Status)

	orderStatus := model.AWAITINGPAYMENT

	_, err = s.stockRepository.Reserve(ctx, items)
	if err != nil {
		logger.Errorw("error reserve in loms", "error", err)
		orderStatus = model.FAILED
	}

	err = s.orderRepository.SetStatus(ctx, *upOrder, orderStatus)
	if err != nil {
		logger.Errorw("error set status in repository", "error", err)
		return upOrder.OrderID, err // Заказ уже записан в статусе new. Так что id можно вернуть.
	}

	s.outboxService.CreateMessage(ctx, upOrder.OrderID, orderStatus)

	if orderStatus == model.FAILED {
		return upOrder.OrderID, model.ErrOutOfStock // Заказ уже записан в статусе new. Так что id можно вернуть.
	}

	return upOrder.OrderID, nil
}
