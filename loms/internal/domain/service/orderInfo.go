package service

import (
	"context"
	"errors"
	"route256/loms/internal/domain/model"
	"sort"
)

func (s *LomsService) OrderInfo(ctx context.Context, orderID int64) (*model.Order, error) {
	if orderID < 1 {

		return nil, model.ErrOrderIDIsNotValid
	}

	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, model.ErrOrderDoesntExist) {
			return nil, model.ErrOrderDoesntExist
		}
		return nil, err
	}

	sort.Slice(order.Items, func(i, j int) bool { return order.Items[i].Sku < order.Items[j].Sku })

	return order, nil
}
