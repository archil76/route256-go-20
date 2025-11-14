package hw1

import (
	"context"
	"math"
	"net/http"

	"route256/tests/app/assert"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestAddSkuWrongUserError(t provider.T) {
	t.Parallel()

	t.Title("Неуспешное добавление SKU в корзину с невалидным userID")

	var (
		sku            = int64(1076963)
		zeroUserID     = int64(0)
		negativeUserID = int64(math.MinInt64)
		ctx            = context.Background()
	)

	t.WithNewParameters(
		"Нулевой userID", zeroUserID,
		"Отрицательный userID", negativeUserID,
		"Основной SKU", sku,
	)

	t.WithNewStep("Добавление SKU в корзину с нулевым userID", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, zeroUserID, sku, 1)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Проверка корзины", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, zeroUserID)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		assert.EmptyCart(sCtx, cart)
	})

	t.WithNewStep("Добавление товара с отрицательным userID", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, negativeUserID, sku, 1)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Проверка корзины", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, negativeUserID)

		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		assert.EmptyCart(sCtx, cart)
	})
}
