package service

import (
	"context"
	"errors"
	"route256/cart/internal/domain/model"
	productrepo "route256/cart/internal/domain/repository/productservicerepository"
	"route256/cart/internal/infra/logger"
)

func (s *CartService) AddItem(ctx context.Context, userID model.UserID, skuID model.Sku, count uint32) (model.Sku, error) {
	ctx, span := s.tracer.Start(ctx, "Service.AddItem")
	defer span.End()

	if userID < 1 || skuID < 1 || count < 1 {
		return 0, ErrFailValidation
	}

	product, err := s.productService.GetProductBySku(ctx, skuID)
	if err != nil {
		logger.Errorw("productService", err)
		if errors.Is(err, productrepo.ErrProductNotFound) {
			return 0, model.ErrProductNotFound
		}
		return 0, err
	}

	if product != nil {
		if product.Sku != skuID {
			return 0, ErrFailValidation
		}
	} else {
		return 0, ErrFailValidation
	}

	countInStock, err := s.lomsService.StockInfo(ctx, skuID)
	if err != nil || countInStock < count {
		logger.Errorw("lomsService", err)
		return 0, model.ErrProductNotFound
	}

	_, err = s.repository.AddItem(ctx, userID, model.Item{Sku: skuID, Count: count})
	if err != nil {
		logger.Errorw("Service", err)
		return 0, err
	}

	return skuID, nil
}
