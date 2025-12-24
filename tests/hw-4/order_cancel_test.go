package hw4

import (
	"context"
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
	"route256/tests/app/domain"
)

func (s *Suite) TestOrderCancel_Success(t provider.T) {
	t.Title("Успешная отмена заказа")

	const (
		userID = 42
		sku    = 139275865
		count  = 2
	)

	var (
		orderID int64
		ctx     = context.Background()
	)

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		var statusCode int
		orderID, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: count},
		})
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)
	})

	t.WithNewStep("Отмена заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderCancel(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Проверка заказа", func(sCtx provider.StepCtx) {
		order, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expectedOrder := &domain.Order{
			Status: domain.OrderStatusCancelled,
			User:   userID,
			Items: []domain.OrderItem{
				{Sku: sku, Count: count},
			},
		}
		assert.Order(sCtx, expectedOrder, order)
	})
}

func (s *Suite) TestOrderCancel_CancelledOrder(t provider.T) {
	t.Title("Отмена уже отмененного заказа")

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

	t.WithNewStep("Отмена заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderCancel(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Проверка заказа", func(sCtx provider.StepCtx) {
		order, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expectedOrder := &domain.Order{
			Status: domain.OrderStatusCancelled,
			User:   userID,
			Items: []domain.OrderItem{
				{Sku: sku, Count: count},
			},
		}
		assert.Order(sCtx, expectedOrder, order)
	})

	t.WithNewStep("Повторная отмена заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderCancel(ctx, sCtx, orderID)
		// Идемпотентность требует успешного ответа при повторной отмене
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Оплата отмененного заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderPay(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})
}

func (s *Suite) TestOrderCancel_PaidOrder(t provider.T) {
	t.Title("Отмена оплаченного заказа")

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

	t.WithNewStep("Отмена оплаченного заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderCancel(ctx, sCtx, orderID)
		// Считаем что отменить оплаченный заказ нельзя, требуется возврат средств, через тех. поддержку
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})
}

func (s *Suite) TestOrderCancel_NotExistingOrder(t provider.T) {
	t.Title("Отмена несуществующего заказа")

	var (
		ctx = context.Background()
	)

	t.WithNewStep("Отмена заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderCancel(ctx, sCtx, 9223372036854775001)
		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
	})
}

func (s *Suite) TestOrderCancel_InvalidRequest(t provider.T) {
	t.Title("Неуспешная отмена из-за невалидного запроса")

	var (
		ctx = context.Background()
	)

	t.Run("Нулевой идентификатор заказа", func(t provider.T) {
		t.WithNewStep("Отмена заказа", func(sCtx provider.StepCtx) {
			statusCode := s.lomsClient.OrderCancel(ctx, sCtx, 0)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Отрицательный идентификатор заказа", func(t provider.T) {
		t.WithNewStep("Отмена заказа", func(sCtx provider.StepCtx) {
			statusCode := s.lomsClient.OrderCancel(ctx, sCtx, -1)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})
}
