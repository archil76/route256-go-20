package hw4

import (
	"context"
	"net/http"
	"route256/tests/app/domain"
	"sync"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
)

func (s *Suite) TestOrderInfo_NotExistingOrder(t provider.T) {
	t.Title("Получение несуществующего заказа")

	var (
		ctx = context.Background()
	)

	t.WithNewStep("Получение заказа", func(sCtx provider.StepCtx) {
		_, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, 9223372036854775001)
		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
	})
}

func (s *Suite) TestOrderInfo_WhileChanging(t provider.T) {
	t.Title("Получение информации о заказе во время его изменения")

	const (
		userID = 42
		sku    = 139275865
		count  = 2
	)

	var (
		orderID int64
		ctx     = context.Background()
		wg      sync.WaitGroup
	)

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		var statusCode int
		orderID, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: count},
		})
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)
	})

	t.WithNewStep("Параллельное получение информации и оплата заказа", func(sCtx provider.StepCtx) {
		var (
			orderInfo  *domain.Order
			statusCode int
		)

		wg.Add(2)
		go func() {
			defer wg.Done()
			orderInfo, statusCode = s.lomsClient.OrderInfo(ctx, sCtx, orderID)
		}()
		go func() {
			defer wg.Done()
			s.lomsClient.OrderPay(ctx, sCtx, orderID)
		}()
		wg.Wait()

		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		sCtx.Require().NotNil(orderInfo)
	})
}

func (s *Suite) TestOrderInfo_ConcurrentRead(t provider.T) {
	t.Title("Параллельное получение информации о заказе")

	const (
		userID = 42
		sku    = 139275865
		count  = 2
	)

	var (
		orderID int64
		ctx     = context.Background()
		wg      sync.WaitGroup
	)

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		var statusCode int
		orderID, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: count},
		})
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)
	})

	t.WithNewStep("Параллельное получение информации о заказе", func(sCtx provider.StepCtx) {
		const concurrentReads = 10
		results := make([]struct {
			order      *domain.Order
			statusCode int
		}, concurrentReads)

		for i := range concurrentReads {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				results[index].order, results[index].statusCode = s.lomsClient.OrderInfo(ctx, sCtx, orderID)
			}(i)
		}
		wg.Wait()

		// Проверяем, что все запросы вернули одинаковую информацию
		for i := 1; i < concurrentReads; i++ {
			assert.StatusCode(sCtx, http.StatusOK, results[i].statusCode)
			sCtx.Require().Equal(results[0].order, results[i].order, "Информация о заказе должна быть консистентной")
		}
	})
}

func (s *Suite) TestOrderInfo_InvalidRequest(t provider.T) {
	t.Title("Неуспешное получение заказа из-за невалидного запроса")

	var (
		ctx = context.Background()
	)

	t.Run("Нулевой идентификатор заказа", func(t provider.T) {
		t.WithNewStep("Получение заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, 0)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Отрицательный идентификатор заказа", func(t provider.T) {
		t.WithNewStep("Получение заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, -1)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})
}
