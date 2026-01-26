package loms

import (
	"context"
	"errors"
	"route256/loms/internal/domain/model"
)

func (s *LomsService) OrderCancel(ctx context.Context, orderID int64) error {
	if orderID < 1 {
		return model.ErrOrderIDIsNotValid
	}

	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, model.ErrOrderDoesntExist) {
			return model.ErrOrderDoesntExist
		}
		return err
	}

	if order.Status == model.CANCELED || order.Status == model.NEWSTATUS {
		return nil // стоки увеличивать не надо
	}

	if order.Status == model.PAYED || order.Status == model.FAILED {
		return model.ErrInvalidOrderStatus
	}

	err = s.stockRepository.ReserveCancel(ctx, order.Items)
	if err != nil {
		return err
	}

	err = s.orderRepository.SetStatus(ctx, *order, model.CANCELED)
	if err != nil {
		return err
	}

	s.outboxService.CreateMessage(ctx, orderID, model.CANCELED)

	return nil
}
