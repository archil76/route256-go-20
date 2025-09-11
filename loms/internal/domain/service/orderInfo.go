package service

import (
	"context"
	"errors"
	"route256/loms/internal/domain/model"
	orderrepo "route256/loms/internal/domain/repository/inmemoryrepository/order"
)

func (s *LomsService) OrderInfo(ctx context.Context, orderID int64) (*model.Order, error) {
	if orderID < 1 {

		return nil, ErrOrderIDIsNotValid
	}

	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, orderrepo.ErrOrderDoesntExist) {
			return nil, ErrOrderDoesntExist
		}
		return nil, err
	}

	return order, nil
}
