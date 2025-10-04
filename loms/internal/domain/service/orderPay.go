package service

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (s *LomsService) OrderPay(ctx context.Context, orderID int64) error {
	if orderID < 1 {
		return model.ErrOrderIDIsNotValid
	}

	order, err := s.OrderInfo(ctx, orderID)
	if err != nil {
		return model.ErrOrderDoesntExist
	}

	if order.Status == model.PAYED {
		return nil // тут разногласия в спеке видимо если оплатили то будет аванс, а вот стоки уменьшать не надо
	}

	if order.Status != model.AWAITINGPAYMENT {
		return model.ErrInvalidOrderStatus
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
