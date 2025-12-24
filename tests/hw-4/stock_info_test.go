package hw4

import (
	"context"
	"net/http"
	"sync"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
	"route256/tests/app/domain"
)

func (s *Suite) TestStocksInfo_Success(t provider.T) {
	t.Title("Успешное получение стоков")

	var (
		ctx = context.Background()
	)

	t.WithNewStep("Получение стоков", func(sCtx provider.StepCtx) {
		_, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, 139275865)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})
}

func (s *Suite) TestStocksInfo_PayOrder(t provider.T) {
	t.Title("Успешное получение стоков после оплаты заказа")

	const (
		userID = 42
		sku    = 139275865
		count  = 2
	)

	var (
		ctx            = context.Background()
		initStockCount uint64
	)

	t.WithNewStep("Получение изначально количества стоков", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		initStockCount = stockCount
	})

	var orderID int64

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		var statusCode int
		orderID, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: count},
		})
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)
	})

	t.WithNewStep("Получение стоков", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.Stocks(sCtx, initStockCount-count, stockCount)
	})

	t.WithNewStep("Оплата заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderPay(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Получение стоков", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.Stocks(sCtx, initStockCount-count, stockCount)
	})
}

func (s *Suite) TestStocksInfo_CancelOrder(t provider.T) {
	t.Title("Успешное получение стоков после отмены заказа")

	const (
		userID = 42
		sku    = 139275865
		count  = 3
	)

	var (
		ctx            = context.Background()
		initStockCount uint64
	)

	t.WithNewStep("Получение стоков", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		initStockCount = stockCount
	})

	var orderID int64

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		var statusCode int
		orderID, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: count},
		})
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)
	})

	t.WithNewStep("Проверка стоков после создания заказа", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.Stocks(sCtx, initStockCount-count, stockCount)
	})

	t.WithNewStep("Отмена заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderCancel(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Проверка стоков после отмены заказа", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.Stocks(sCtx, initStockCount, stockCount)
	})
}

func (s *Suite) TestStocksInfo_NonExistentSKU(t provider.T) {
	t.Title("Получение информации о несуществующем SKU")

	const (
		nonExistentSKU = 999999999
	)

	var (
		ctx = context.Background()
	)

	t.WithNewStep("Получение информации о несуществующем SKU", func(sCtx provider.StepCtx) {
		_, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, nonExistentSKU)
		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
	})
}

func (s *Suite) TestStocksInfo_InvalidRequest(t provider.T) {
	t.Title("Неуспешное получение стоков из-за невалидной запроса")

	var (
		ctx = context.Background()
	)

	t.Run("Нулевой идентификатор SKU", func(t provider.T) {
		t.WithNewStep("Получение стоков SKU", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, 0)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Отрицательный идентификатор SKU", func(t provider.T) {
		t.WithNewStep("Получение стоков SKU", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, -1)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})
}

func (s *Suite) TestStocksInfo_AddCartItem(t provider.T) {
	t.Title("Неизменность стоков при добавлении товара в корзину")

	const (
		userID = 42
		sku    = 2956315
		count  = 3
	)

	var (
		ctx            = context.Background()
		initStockCount uint64
	)

	t.WithNewStep("Получение стоков", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		initStockCount = stockCount
	})

	t.WithNewStep("Добавление товаров в корзину", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku, count)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Проверка стоков после добавления товаров в корзину", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.Stocks(sCtx, initStockCount, stockCount)
	})
}

func (s *Suite) TestStocksInfo_ConcurrentRead(t provider.T) {
	t.Title("Параллельное получение информации о стоках")

	const (
		sku = 139275865
	)

	var (
		ctx = context.Background()
		wg  sync.WaitGroup
	)

	t.WithNewStep("Параллельное получение информации о стоках", func(sCtx provider.StepCtx) {
		const concurrentReads = 10
		results := make([]struct {
			count      uint64
			statusCode int
		}, concurrentReads)

		for i := range concurrentReads {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				results[index].count, results[index].statusCode = s.lomsClient.StocksInfo(ctx, sCtx, sku)
			}(i)
		}
		wg.Wait()

		// Проверяем, что все запросы вернули одинаковое количество стоков
		for i := 1; i < concurrentReads; i++ {
			assert.StatusCode(sCtx, http.StatusOK, results[i].statusCode)
			sCtx.Require().Equal(results[0].count, results[i].count, "Количество стоков должно быть консистентным")
		}
	})
}
