package service

import (
	"context"
	"route256/cart/internal/domain/model"
	"sort"
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

	err = s.fillReportCart(ctx, cart, &reportCart)
	if err != nil {
		return nil, err
	}
	return &reportCart, nil
}

func (s *CartService) fillReportCart(ctx context.Context, cart *model.Cart, reportCart *model.ReportCart) error {
	skus := make([]model.Sku, 0, len(cart.Items))
	for sku := range cart.Items {
		skus = append(skus, sku)
	}

	products, err := s.productService.GetProductsBySkus(ctx, skus)
	if err != nil {
		return err
	}

	totalPrice := uint32(0)
	for _, product := range products {
		count := cart.Items[product.Sku]
		reportCart.Items = append(reportCart.Items, model.ItemInCart{
			SKU:   product.Sku,
			Count: count,
			Name:  product.Name,
			Price: product.Price,
		})

		totalPrice += product.Price * count
	}

	reportCart.TotalPrice += totalPrice

	sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })

	return nil
}
