package hw3

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
	"route256/tests/app/domain"
)

func (s *Suite) TestCheckout_Success(t provider.T) {
	t.Title("Успешное оформление заказа")

	const (
		sku1   = 139275865
		sku2   = 1076963
		count1 = 1
		count2 = 2
	)

	var (
		ctx     = context.Background()
		userID  = int64(rand.Int())
		orderID int64
	)

	t.WithNewStep("Проверка корзины", func(sCtx provider.StepCtx) {
		_, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
	})

	t.WithNewStep("Добавление товаров в корзину", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku1, count1)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		statusCode = s.cartClient.AddItem(ctx, sCtx, userID, sku2, count2)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Проверка корзины", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expectedCart := &domain.Cart{
			Items: []domain.CartItem{
				{
					SKU:   sku2,
					Count: count2,
					Name:  "Теория нравственных чувств | Смит Адам",
					Price: 3379,
				},
				{
					SKU:   sku1,
					Count: count1,
					Name:  "Платье oodji Collection",
					Price: 634,
				},
			},
			TotalPrice: 7392,
		}
		assert.Cart(sCtx, expectedCart, cart)
	})

	t.WithNewStep("Оформление заказа", func(sCtx provider.StepCtx) {
		var statusCode int
		orderID, statusCode = s.cartClient.Checkout(ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)
	})

	t.WithNewStep("Проверка заказа", func(sCtx provider.StepCtx) {
		order, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expectedOrder := &domain.Order{
			Status: domain.OrderStatusAwaitingPayment,
			User:   userID,
			Items: []domain.OrderItem{
				{Sku: sku2, Count: count2},
				{Sku: sku1, Count: count1},
			},
		}
		assert.Order(sCtx, expectedOrder, order)
	})

	t.WithNewStep("Проверка корзины", func(sCtx provider.StepCtx) {
		_, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
	})
}

func (s *Suite) TestCheckout_EmptyCart(t provider.T) {
	t.Title("Неуспешное оформление заказа из-за отсутствия товаров в корзине")

	const userID int64 = 13377331 // Не должно пересекаться в других тестах.

	var (
		ctx = context.Background()
	)

	t.WithNewStep("Оформление заказа", func(sCtx provider.StepCtx) {
		_, statusCode := s.cartClient.Checkout(ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
	})
}

func (s *Suite) TestCheckout_InvalidRequest(t provider.T) {
	t.Title("Неуспешная оформление заказа из-за невалидного запроса")

	var (
		ctx = context.Background()
	)

	t.Run("Нулевой идентификатор пользователя", func(t provider.T) {
		t.WithNewStep("Оформление заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.cartClient.Checkout(ctx, sCtx, 0)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Отрицательный идентификатор пользователя", func(t provider.T) {
		t.WithNewStep("Оформление заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.cartClient.Checkout(ctx, sCtx, -1)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})
}
