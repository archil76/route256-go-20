package hw1

import (
	"context"
	"math/rand"
	"net/http"

	"route256/tests/app/assert"
	"route256/tests/app/domain"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestDeleteSkuSuccess(t provider.T) {
	t.Parallel()

	t.Title("Успешное удаление SKU из корзины")

	var (
		sku    = int64(1076963)
		sku2   = int64(1148162)
		count  = int64(2)
		count2 = int64(3)
		userID = int64(rand.Int())
		ctx    = context.Background()
	)

	t.WithNewParameters(
		"userID", userID,
		"Основной SKU", sku,
		"Количество экземпляров основного SKU", count,
		"Дополнительный SKU", sku2,
		"Количество экземпляров дополнительного SKU", count2,
	)

	t.WithTestSetup(func(t provider.T) {

		t.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
			statusCode := s.cartClient.DeleteCart(ctx, sCtx, userID)

			assert.StatusCode(sCtx, http.StatusNoContent, statusCode)
		})

		t.WithNewStep("Проверка что корзина пуста", func(sCtx provider.StepCtx) {
			cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

			assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
			assert.EmptyCart(sCtx, cart)
		})

		t.WithNewStep("Наполнение корзины", func(sCtx provider.StepCtx) {
			statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku, count)
			assert.StatusCode(sCtx, http.StatusOK, statusCode)

			statusCode = s.cartClient.AddItem(ctx, sCtx, userID, sku2, count2)
			assert.StatusCode(sCtx, http.StatusOK, statusCode)

		})

		t.WithNewStep("Проверка корзины", func(sCtx provider.StepCtx) {
			cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

			assert.StatusCode(sCtx, http.StatusOK, statusCode)
			assert.NotEmptyCart(sCtx, cart)
			sCtx.Assert().Len(cart.Items, 2, "Ожидается два товара в корзине")
		})
	})

	t.WithNewStep("Удаление SKU из корзины", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.DeleteItem(ctx, sCtx, userID, sku)

		assert.StatusCode(sCtx, http.StatusNoContent, statusCode)
	})

	t.WithNewStep("Проверка корзины после удаления", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expectedCart := &domain.Cart{
			Items: []domain.CartItem{
				{
					SKU:   uint32(sku2),
					Count: count2,
					Name:  cart.Items[0].Name,
					Price: cart.Items[0].Price,
				},
			},
			TotalPrice: cart.TotalPrice,
		}
		assert.Cart(sCtx, expectedCart, cart)
		assert.NotEmptyCart(sCtx, cart)
	})
}
