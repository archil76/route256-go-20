package service

import (
	"context"
	"errors"
	"route256/cart/internal/domain/model"
)

var (
	ErrInvalidSKU     = errors.New("invalid sku")
	ErrFailValidation = errors.New("fail validation")
	//ErrCartIsEmpty    = errors.New("cart is empty")
)

type CartsRepository interface {
	AddItem(_ context.Context, userID model.UserID, item model.Item) (*model.Item, error)
	DeleteItem(_ context.Context, userID model.UserID, item model.Item) (*model.Item, error)
	GetCart(_ context.Context, userID model.UserID) (*model.Cart, error)
	DeleteItems(_ context.Context, userID model.UserID) (model.UserID, error)
}

type ProductService interface {
	GetProductBySku(ctx context.Context, sku model.Sku) (*model.Product, error)
}

type CartService struct {
	repository     CartsRepository
	productService ProductService
}

func NewCartsService(repository CartsRepository, productService ProductService) *CartService {
	return &CartService{repository: repository, productService: productService}
}

func (s *CartService) AddItem(ctx context.Context, userID model.UserID, skuID model.Sku, count uint16) (model.Sku, error) {

	if userID < 1 || skuID < 1 || count < 1 {
		return 0, ErrFailValidation
	}

	if _, err := s.productService.GetProductBySku(ctx, skuID); err != nil {
		return 0, ErrInvalidSKU
	}

	_, err := s.repository.AddItem(ctx, userID, model.Item{Sku: skuID, Count: count})
	if err != nil {
		return 0, err
	}

	return skuID, nil
}

func (s *CartService) DeleteItem(ctx context.Context, userID model.UserID, skuID model.Sku) (model.Sku, error) {

	if userID < 1 || skuID < 1 {
		return 0, ErrFailValidation
	}

	_, err := s.repository.DeleteItem(ctx, userID, model.Item{Sku: skuID, Count: 0})
	if err != nil {
		return 0, err
	}

	return skuID, nil

}

func (s *CartService) GetItemsByUserID(ctx context.Context, userID model.UserID) (*model.ReportCart, error) {

	if userID < 1 {
		return nil, ErrFailValidation
	}

	cart, err := s.repository.GetCart(ctx, userID)

	if err != nil {
		return nil, err
	}

	//if cart.Items == nil || len(cart.Items) == 0 {
	// // can't be reach
	//	return nil, ErrCartIsEmpty
	//}

	reportCart := model.ReportCart{
		UserID:     userID,
		Items:      map[model.Sku]model.ItemInСart{},
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

		reportCart.Items[sku] = model.ItemInСart{
			SKU:   sku,
			Count: count,
			Name:  name,
			Price: price,
		}

		reportCart.TotalPrice += price * uint32(count)

	}

	return &reportCart, nil
}

func (s *CartService) DeleteItemByUserID(ctx context.Context, userID model.UserID) (model.UserID, error) {

	if userID < 1 {
		return 0, ErrFailValidation
	}
	_, err := s.repository.DeleteItems(ctx, userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}
