package service

import (
	"context"
	"route256/cart/internal/domain/model"
	"route256/cart/internal/infra/errgroup"
	"sort"
	"sync"
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

// Зависимость от errgroup через интерфейс
// Отвязать зависимость от ReportCart. Т.е возвращать массив продактов
func (s *CartService) fillReportCart(ctx context.Context, cart *model.Cart, reportCart *model.ReportCart) error {
	mu := &sync.Mutex{}
	group, ctx := errgroup.WithContext(ctx, 10, len(cart.Items))
	group.RunWorker()

	totalPrice := uint32(0)
	for sku, count := range cart.Items {
		group.Go(func() error {
			name := ""
			price := uint32(0)

			product, err := s.productService.GetProductBySku(ctx, sku)
			if err != nil {
				ctx.Done()
				return err
			}
			name = product.Name
			price = product.Price

			mu.Lock()
			defer mu.Unlock()
			reportCart.Items = append(reportCart.Items, model.ItemInCart{
				SKU:   sku,
				Count: count,
				Name:  name,
				Price: price,
			})

			totalPrice += price * count
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return err
	}

	reportCart.TotalPrice += totalPrice

	sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })

	return nil
}
