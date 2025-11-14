package hw1

import (
	"context"
	"math"
	"math/rand"
	"net/http"

	"route256/tests/app/assert"
	"route256/tests/app/domain"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestDeleteNonExistSkuSuccess(t provider.T) {
	t.Parallel()

	t.Title("Удаление несуществующего SKU из корзины")

	var (
		sku         = int64(1076963)
		nonExistSku = int64(999999)
		userID      = rand.Int63()
		ctx         = context.Background()
	)

	t.WithNewParameters(
		"userID", userID,
		"Существующий SKU", sku,
		"Несуществующий SKU", nonExistSku,
	)

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
			statusCode := s.cartClient.DeleteCart(ctx, sCtx, userID)

			assert.StatusCode(sCtx, http.StatusNoContent, statusCode)
		})

		t.WithNewStep("Добавление SKU в корзину", func(sCtx provider.StepCtx) {
			statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku, 1)

			assert.StatusCode(sCtx, http.StatusOK, statusCode)
		})
	})

	t.WithNewStep("Удаление несуществующего SKU из корзины", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.DeleteItem(ctx, sCtx, userID, nonExistSku)

		assert.StatusCode(sCtx, http.StatusNoContent, statusCode)
	})

	t.WithNewStep("Проверка корзины после попытки удаления", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expectedCart := &domain.Cart{
			Items: []domain.CartItem{
				{
					SKU:   uint32(sku),
					Count: 1,
					Name:  cart.Items[0].Name,
					Price: cart.Items[0].Price,
				},
			},
			TotalPrice: cart.TotalPrice,
		}
		assert.Cart(sCtx, expectedCart, cart)
	})
}

func (s *Suite) TestDeleteInvalidSkuError(t provider.T) {
	t.Parallel()

	t.Title("Неуспешное удаление SKU с невалидными параметрами")

	var (
		userID      = rand.Int63()
		zeroSku     = int64(0)
		negativeSku = int64(math.MinInt64)
		ctx         = context.Background()
	)

	t.WithNewParameters(
		"userID", userID,
		"Нулевой SKU", zeroSku,
		"Отрицательный SKU", negativeSku,
	)

	t.WithNewStep("Удаление SKU с нулевым значением", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.DeleteItem(ctx, sCtx, userID, zeroSku)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Удаление SKU с отрицательным значением", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.DeleteItem(ctx, sCtx, userID, negativeSku)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})
}

func (s *Suite) TestDeleteSkuWrongUserError(t provider.T) {
	t.Parallel()

	t.Title("Неуспешное удаление SKU с невалидным userID")

	var (
		sku            = int64(1076963)
		zeroUserID     = int64(0)
		negativeUserID = int64(math.MinInt64)
		ctx            = context.Background()
	)

	t.WithNewParameters(
		"Нулевой userID", zeroUserID,
		"Отрицательный userID", negativeUserID,
		"SKU", sku,
	)

	t.WithNewStep("Удаление SKU с нулевым userID", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.DeleteItem(ctx, sCtx, zeroUserID, sku)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Удаление SKU с отрицательным userID", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.DeleteItem(ctx, sCtx, negativeUserID, sku)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})
}
