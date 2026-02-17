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
	var upOrder *model.Order
	err := s.pooler.InTx(ctx, func(ctx context.Context) error {
		var err error
		upOrder, err = s.orderRepository.Create(ctx, order)
		if err != nil {
			return err
		}

		s.outboxService.CreateMessage(ctx, upOrder.OrderID, upOrder.Status)

		return nil
	})

	if err != nil {
		return 0, err
	}

	err = s.pooler.InTx(ctx, func(ctx context.Context) error {
		var err error
		orderStatus := model.AWAITINGPAYMENT

		_, err = s.stockRepository.Reserve(ctx, items)
		if err != nil {
			logger.Errorw("error reserve in loms", "error", err)
			orderStatus = model.FAILED
		}
		upOrder.Status = orderStatus
		err = s.orderRepository.SetStatus(ctx, *upOrder, orderStatus)
		if err != nil {
			logger.Errorw("error set status in repository", "error", err)
			return err // Заказ уже записан в статусе new. Так что id можно вернуть.
		}

		s.outboxService.CreateMessage(ctx, upOrder.OrderID, orderStatus)

		return nil
	})

	if err != nil {
		return upOrder.OrderID, err
	}

	if upOrder.Status == model.FAILED {
		return upOrder.OrderID, model.ErrOutOfStock // Заказ уже записан в статусе new. Так что id можно вернуть.
	}

	return upOrder.OrderID, nil
}
