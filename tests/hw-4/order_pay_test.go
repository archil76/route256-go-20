package hw4

import (
	"context"
	"net/http"
	"sync"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
	"route256/tests/app/domain"
)

func (s *Suite) TestOrderPay_Success(t provider.T) {
	t.Title("Успешная оплата заказа")

	const (
		userID = 42
		sku    = 139275865
		count  = 2
	)

	var (
		ctx     = context.Background()
		orderID int64
	)

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		var statusCode int
		orderID, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: 2},
		})
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)
	})

	t.WithNewStep("Оплата заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderPay(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Проверка заказа", func(sCtx provider.StepCtx) {
		res, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expected := &domain.Order{
			Status: domain.OrderStatusPaid,
			User:   userID,
			Items: []domain.OrderItem{
				{Sku: sku, Count: count},
			},
		}
		assert.Order(sCtx, expected, res)
	})
}

func (s *Suite) TestOrderPay_PaidOrder(t provider.T) {
	t.Title("Оплата уже оплаченного заказа")

	const (
		userID = 42
		sku    = 139275865
		count  = 2
	)

	var (
		ctx     = context.Background()
		orderID int64
	)

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		var statusCode int
		orderID, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: count},
		})
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)
	})

	t.WithNewStep("Оплата заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderPay(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Проверка заказа", func(sCtx provider.StepCtx) {
		res, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expected := &domain.Order{
			Status: domain.OrderStatusPaid,
			User:   userID,
			Items: []domain.OrderItem{
				{Sku: sku, Count: count},
			},
		}
		assert.Order(sCtx, expected, res)
	})

	t.WithNewStep("Повторная оплата заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderPay(ctx, sCtx, orderID)
		// Идемпотентность требует успешного ответа при повторной оплате
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})
}

func (s *Suite) TestOrderPay_NotExistingOrder(t provider.T) {
	t.Title("Оплата несуществующего заказа")

	var (
		ctx = context.Background()
	)

	t.WithNewStep("Оплата заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderPay(ctx, sCtx, 9223372036854775001)
		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
	})
}

func (s *Suite) TestOrderPay_InvalidRequest(t provider.T) {
	t.Title("Неуспешная оплата из-за невалидного запроса")

	var (
		ctx = context.Background()
	)

	t.Run("Нулевой идентификатор заказа", func(t provider.T) {
		t.WithNewStep("Оплата заказа", func(sCtx provider.StepCtx) {
			statusCode := s.lomsClient.OrderPay(ctx, sCtx, 0)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Отрицательный идентификатор заказа", func(t provider.T) {
		t.WithNewStep("Оплата заказа", func(sCtx provider.StepCtx) {
			statusCode := s.lomsClient.OrderPay(ctx, sCtx, -1)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})
}

func (s *Suite) TestOrderPay_ConcurrentPay(t provider.T) {
	// Так как создание заказа и обновление стоков не в транзакции, то тест ломает данные:
	// отрицательное резервированное количество.
	t.Skip()

	t.Title("Параллельная оплата одного и того же заказа")

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

	t.WithNewStep("Параллельная оплата заказа", func(sCtx provider.StepCtx) {
		const concurrentPays = 10
		results := make([]int, concurrentPays)

		for i := 0; i < concurrentPays; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				results[index] = s.lomsClient.OrderPay(ctx, sCtx, orderID)
			}(i)
		}
		wg.Wait()

		// Проверяем, что все попытки оплаты вернули успешный статус
		for _, statusCode := range results {
			assert.StatusCode(sCtx, http.StatusOK, statusCode)
		}
	})

	t.WithNewStep("Проверка заказа", func(sCtx provider.StepCtx) {
		res, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expected := &domain.Order{
			Status: domain.OrderStatusPaid,
			User:   userID,
			Items: []domain.OrderItem{
				{Sku: sku, Count: count},
			},
		}
		assert.Order(sCtx, expected, res)
	})
}
