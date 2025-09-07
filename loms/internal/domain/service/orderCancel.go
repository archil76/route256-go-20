package service

import (
	"context"
	"errors"
	"route256/loms/internal/domain/model"
	orderrepo "route256/loms/internal/domain/repository/inmemoryrepository/order"
)

func (s *LomsService) OrderCancel(ctx context.Context, orderId int64) error {
	if orderId < 1 {
		return ErrOrderIDIsNotValid
	}

	order, err := s.orderRepository.GetByID(ctx, orderId)
	if err != nil {
		if errors.Is(err, orderrepo.ErrOrderDoesntExist) {
			return ErrOrderDoesntExist
		}
		return err
	}

	if order.Status == model.CANCELED || order.Status == model.NEW_STATUS {
		return nil // стоки увеличивать не надо
	}

	if order.Status == model.PAYED || order.Status == model.FAILED {
		return ErrInvalidOrderStatus
	}

	err = s.stockRepository.ReserveCancel(ctx, order.Items)
	if err != nil {
		return err
	}

	err = s.orderRepository.SetStatus(ctx, *order, model.CANCELED)
	if err != nil {
		return err
	}

	return nil
}
