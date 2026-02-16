//go:build e2e_test

package hw3e2e

import (
	"net/http"
	"route256/tests/app/assert"
	"route256/tests/app/domain"
	"sync"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"golang.org/x/net/context"
)

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
