package service

import (
	"context"

	"route256/cart/internal/domain/model"
)

func (s *CartService) GetItemsByUserID(ctx context.Context, userID model.UserID) (*model.ReportCart, error) {

	if userID < 1 {
		return nil, ErrFailValidation
	}

	cart, err := s.repository.GetCart(ctx, userID)

	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, ErrCartIsEmpty
	}

	reportCart := model.ReportCart{
		Items:      []model.ItemInСart{},
		TotalPrice: 0,
	}

	for sku, count := range cart.Items {
		name := ""
		price := uint32(0)

		itemInfo, err := s.productService.GetProductBySku(ctx, sku)
		if err == nil {
			name = itemInfo.Name
			price = itemInfo.Price
		}

		reportCart.Items = append(reportCart.Items, model.ItemInСart{
			SKU:   sku,
			Count: count,
			Name:  name,
			Price: price,
		})

		reportCart.TotalPrice += price * uint32(count)

	}

	return &reportCart, nil
}
