package service

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (s *LomsService) OrderPay(ctx context.Context, orderId int64) error {
	if orderId < 1 {
		return ErrOrderIDIsNotValid
	}

	order, err := s.OrderInfo(ctx, orderId)
	if err != nil {
		return ErrOrderDoesntExist
	}

	if order.Status == model.PAYED {
		return nil // тут разногласия в спеке видимо если оплатили то будет аванс, а вот стоки уменьшать не надо
	}

	if order.Status != model.AWAITING_PAYMENT {
		return ErrInvalidOrderStatus
	}

	err = s.stockRepository.ReserveRemove(ctx, order.Items)
	if err != nil {
		return err
	}

	err = s.orderRepository.SetStatus(ctx, *order, model.PAYED)
	if err != nil {
		return err
	}

	return nil
}
