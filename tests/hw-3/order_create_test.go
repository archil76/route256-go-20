package hw3

import (
	"net/http"
	"route256/tests/app/assert"
	"route256/tests/app/domain"
	"sync"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"golang.org/x/net/context"
)

func (s *Suite) TestOrderCreate_Success(t provider.T) {
	t.Title("Успешное создание заказа")

	const (
		userID = 42
		sku1   = 139275865
		sku2   = 1076963
		count1 = 2
		count2 = 1
	)

	var (
		ctx     = context.Background()
		orderID int64
	)

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		var statusCode int
		orderID, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku1, Count: count1},
			{Sku: sku2, Count: count2},
		})
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)
	})

	t.WithNewStep("Проверка заказа", func(sCtx provider.StepCtx) {
		res, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expected := &domain.Order{
			User:   userID,
			Status: domain.OrderStatusAwaitingPayment,
			Items: []domain.OrderItem{
				{Sku: sku2, Count: count2},
				{Sku: sku1, Count: count1},
			},
		}

		assert.Order(sCtx, expected, res)
	})
}

func (s *Suite) TestOrderCreate_NoStockInfo(t provider.T) {
	t.Title("Неуспешное создание заказ из-за отсутствия информации о стоках товара")

	const (
		userID = 42
		sku    = 404
	)

	var (
		ctx = context.Background()
	)

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		_, statusCode := s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: 1},
		})
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})
}

func (s *Suite) TestOrderCreate_InvalidRequest(t provider.T) {
	t.Title("Неуспешное создание заказа из-за невалидного запроса")

	var (
		ctx = context.Background()
	)

	t.Run("Нулевой идентификатор пользователя", func(t provider.T) {
		t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderCreate(ctx, sCtx, 0, []domain.OrderItem{{Sku: 1, Count: 2}})
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Отрицательный идентификатор пользователя", func(t provider.T) {
		t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderCreate(ctx, sCtx, -1, []domain.OrderItem{{Sku: 1, Count: 2}})
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Нет товаров в заказе", func(t provider.T) {
		t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderCreate(ctx, sCtx, 42, nil)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Нулевой идентификатор для SKU", func(t provider.T) {
		t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderCreate(ctx, sCtx, 42, []domain.OrderItem{
				{Sku: 0, Count: 1},
			})
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Отрицательный идентификатор для SKU", func(t provider.T) {
		t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderCreate(ctx, sCtx, 42, []domain.OrderItem{
				{Sku: -1, Count: 1},
			})
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Нулевое количество товара", func(t provider.T) {
		t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderCreate(ctx, sCtx, 42, []domain.OrderItem{
				{Sku: 1, Count: 0},
			})
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Отрицательное количество товара", func(t provider.T) {
		t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderCreate(ctx, sCtx, 42, []domain.OrderItem{
				{Sku: 1, Count: -1},
			})
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})
}

func (s *Suite) TestOrderCreate_NotEnoughStocks(t provider.T) {
	t.Title("Неуспешное создание заказа из-за недостаточных стоков")

	const (
		sku1 = 1076963
		sku2 = 135937324
	)

	var (
		statusCode int

		skuStock1 uint64
		skuStock2 uint64
		ctx       = context.Background()
	)

	t.WithNewStep("Получение стоков по товарам", func(sCtx provider.StepCtx) {
		skuStock1, statusCode = s.lomsClient.StocksInfo(ctx, sCtx, sku1)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		skuStock2, statusCode = s.lomsClient.StocksInfo(ctx, sCtx, sku2)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		_, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, 42, []domain.OrderItem{
			{Sku: sku1, Count: int64(skuStock1 - 1)},
			{Sku: sku2, Count: int64(skuStock2 + 1)},
		})
		// https://github.com/grpc-ecosystem/grpc-gateway/blob/main/runtime/errors.go#L58
		// grpc-gateway конвертирует codes.FailedPrecondition в http.StatusBadRequest
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Получение стоков по товарам после попытки создания заказа", func(sCtx provider.StepCtx) {
		newSKUStock1, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku1)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.Stocks(sCtx, skuStock1, newSKUStock1)

		newSKUStock2, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku2)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.Stocks(sCtx, skuStock2, newSKUStock2)
	})
}
func (s *Suite) TestOrderCreate_SuccessAsync(t provider.T) {
	t.Title("Успешное создание заказов (async)")

	const (
		sku   = 1625903
		count = 1

		ordersCount = 100
	)

	var (
		userIDs = []int64{42, 43, 44}
		orders  = make([]struct {
			userID  int64
			orderID int64
		}, ordersCount)

		initStocksCount uint64

		ctx = context.Background()
	)

	t.WithNewStep("Получение изначальных стоков", func(sCtx provider.StepCtx) {
		stockCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		initStocksCount = stockCount
	})

	t.WithNewStep("Создание заказов", func(sCtx provider.StepCtx) {
		var wg sync.WaitGroup

		for i := range ordersCount {
			wg.Add(1)
			go func() {
				defer wg.Done()

				userID := userIDs[i%len(userIDs)]
				orderID, statusCode := s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
					{Sku: sku, Count: count},
				})
				assert.StatusCode(sCtx, http.StatusOK, statusCode)
				assert.OrderID(sCtx, orderID)

				orders[i].userID = userID
				orders[i].orderID = orderID
			}()
		}
		wg.Wait()
	})

	t.WithNewStep("Проверка заказов", func(sCtx provider.StepCtx) {
		for _, order := range orders {
			res, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, order.orderID)
			assert.StatusCode(sCtx, http.StatusOK, statusCode)

			expected := &domain.Order{
				Status: domain.OrderStatusAwaitingPayment,
				User:   order.userID,
				Items: []domain.OrderItem{
					{Sku: sku, Count: count},
				},
			}
			assert.Order(sCtx, expected, res)
		}
	})

	t.WithNewStep("Проверка стоков", func(sCtx provider.StepCtx) {
		stocksCount, statusCode := s.lomsClient.StocksInfo(ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		//**
		sCtx.Assert().Equal(int(initStocksCount-ordersCount*count), int(stocksCount))
		//sCtx.Assert().Equal(0, int(ordersCount), "ordersCount")
		//sCtx.Assert().Equal(0, int(count), "count")
		//sCtx.Assert().Equal(int(initStocksCount-ordersCount*count), int(stocksCount), "stocksCount")
		//**

		assert.Stocks(sCtx, initStocksCount-ordersCount*count, stocksCount)
	})
}
