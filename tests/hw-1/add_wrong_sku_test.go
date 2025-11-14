package hw1

import (
	"context"
	"math"
	"math/rand"
	"net/http"

	"route256/tests/app/assert"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestWrongAddSkuSuccess(t provider.T) {
	t.Parallel()

	t.Title("Неуспешное добавление SKU с невалидными параметрами")

	var (
		wrongSku      = uint32(1076963000)
		userID        = int64(rand.Int())
		zeroSku       = int64(0)
		maxAllowedSku = int64(math.MaxInt64)
		negativeSku   = int64(math.MinInt64)
		validSku      = int64(1076963)
		zeroCount     = int64(0)
		negativeCount = int64(-999)
		ctx           = context.Background()
	)

	t.WithNewParameters(
		"userID", userID,
		"Несуществующий SKU", wrongSku,
		"Нулевой SKU", zeroSku,
		"Отрицательный SKU", negativeSku,
		"Невалидный SKU (MaxInt64)", maxAllowedSku,
		"Валидный SKU", validSku,
		"Нулевое количество", zeroCount,
		"Отрицательное количество", negativeCount,
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
	})

	t.WithNewStep("Добавление несуществующего SKU в корзину", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, int64(wrongSku), 1)

		assert.StatusCode(sCtx, http.StatusPreconditionFailed, statusCode)
	})

	t.WithNewStep("Проверка корзины после попытки добавления", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
		assert.EmptyCart(sCtx, cart)
	})

	t.WithNewStep("Добавление SKU с нулевым значением", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, zeroSku, 1)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Добавление SKU с отрицательным значением", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, negativeSku, 1)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Добавление SKU с нулевым количеством", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, validSku, zeroCount)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Добавление SKU с отрицательным количеством", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, validSku, negativeSku)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Попытка добавления SKU с очень большим значением (MaxInt64)", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, maxAllowedSku, 1)
		assert.StatusCode(sCtx, http.StatusPreconditionFailed, statusCode)
	})
}
