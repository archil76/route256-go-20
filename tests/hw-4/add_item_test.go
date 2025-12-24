package hw4

import (
	"context"
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
)

func (s *Suite) TestAddItem_Success(t provider.T) {
	t.Title("Добавление товара с количеством равным стоку")

	const (
		userID = 42
		sku    = 135717466
	)

	var (
		initStockCount uint64
		ctx            = context.Background()
	)

	t.WithNewStep("Получение стоков", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		initStockCount = stockCount
	})

	t.WithNewStep("Добавление товара в корзину", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku, int64(initStockCount))
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})
}

func (s *Suite) TestAddItem_OOS(t provider.T) {
	t.Skip()

	t.Title("Невозможность добавления товара из-за проблемы OOS (out of stock)")

	const (
		userID = 42
		sku    = 1148162
	)

	var (
		initStockCount uint64
		ctx            = context.Background()
	)

	t.WithNewStep("Получение стоков", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		initStockCount = stockCount
	})

	t.WithNewStep("Добавление товара в корзину", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku, int64(initStockCount+1))
		assert.StatusCode(sCtx, http.StatusPreconditionFailed, statusCode)
	})
}
