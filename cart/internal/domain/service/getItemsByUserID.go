package service

import (
	"context"
	"sort"

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
		Items:      []model.ItemInCart{},
		TotalPrice: 0,
	}

	totalPrice := uint32(0)
	for sku, count := range cart.Items {
		name := ""
		price := uint32(0)

		product, err := s.productService.GetProductBySku(ctx, sku)
		if err == nil {
			name = product.Name
			price = product.Price
		}

		reportCart.Items = append(reportCart.Items, model.ItemInCart{
			SKU:   sku,
			Count: count,
			Name:  name,
			Price: price,
		})

		totalPrice += price * count
	}

	reportCart.TotalPrice += totalPrice

	sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })

	return &reportCart, nil
}
