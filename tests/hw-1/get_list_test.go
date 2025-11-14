package hw1

import (
	"context"
	"math"
	"math/rand"
	"net/http"

	"route256/tests/app/assert"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestGetListSuccess(t provider.T) {
	t.Parallel()

	t.Title("Успешное получение содержимого корзины")

	var (
		userID         = rand.Int63()
		skuList        = []int64{1076963, 1148162}
		zeroUserID     = int64(0)
		negativeUserID = int64(math.MinInt64)
		ctx            = context.Background()
	)

	t.WithNewParameters(
		"userID", userID,
		"Список SKU", skuList,
		"Нулевой userID", zeroUserID,
		"Отрицательный userID", negativeUserID,
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

		t.WithNewStep("Наполнение корзины SKU", func(sCtx provider.StepCtx) {
			for _, sku := range skuList {
				statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku, 1)

				assert.StatusCode(sCtx, http.StatusOK, statusCode)
			}
		})
	})

	t.WithNewStep("Получение содержимого корзины", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.NotEmptyCart(sCtx, cart)
		sCtx.Assert().Len(cart.Items, len(skuList), "Ожидается соответствие количества товаров в корзине")
	})

	t.WithNewStep("Удаление корзины", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.DeleteCart(ctx, sCtx, userID)

		assert.StatusCode(sCtx, http.StatusNoContent, statusCode)
	})

	t.WithNewStep("Попытка получения несуществующей корзины", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
		assert.EmptyCart(sCtx, cart)
	})

	t.WithNewStep("Получение корзины с нулевым userID", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, zeroUserID)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		assert.EmptyCart(sCtx, cart)
	})

	t.WithNewStep("Получение корзины с отрицательным userID", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, negativeUserID)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		assert.EmptyCart(sCtx, cart)
	})
}
