package hw7

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ozontech/allure-go/pkg/allure"

	"route256/tests/app/assert"
	"route256/tests/app/domain"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestOrderCancelSuccess(t provider.T) {
	t.Title("Успешная отмена заказа")

	const (
		userID = 42
		sku    = 139275865
		count  = 2
	)

	var (
		orderID    int64
		orderEvent domain.OrderEvent
		ctx        = context.Background()
	)

	t.WithNewParameters(
		"userID", userID,
		"sku", sku,
		"count", count,
	)

	t.WithNewStep("Создание заказа", func(sCtx provider.StepCtx) {
		var statusCode int
		orderID, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: count},
		})
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)

		msg := s.kafkaConsumer.ConsumeSingle(ctx)
		sCtx.Require().NotNil(msg, "Ожидается получение события о создании заказа")
		sCtx.WithNewAttachment("new order message body", allure.JSON, msg.Value)

		err := json.Unmarshal(msg.Value, &orderEvent)
		sCtx.Require().NoError(err, "Ожидается успешная десериализация события")
		assert.OrderEvent(sCtx, domain.OrderEvent{OrderID: orderID, Status: domain.OrderStatusNew}, orderEvent)

		msg = s.kafkaConsumer.ConsumeSingle(ctx)
		sCtx.Require().NotNil(msg, "Ожидается получение события о переходе заказа в статус awaiting payment")
		sCtx.WithNewAttachment("awaiting payment message body", allure.JSON, msg.Value)

		err = json.Unmarshal(msg.Value, &orderEvent)
		sCtx.Require().NoError(err, "Ожидается успешная десериализация события")
		assert.OrderEvent(sCtx, domain.OrderEvent{OrderID: orderID, Status: domain.OrderStatusAwaitingPayment}, orderEvent)
	})

	t.WithNewStep("Отмена заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderCancel(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		msg := s.kafkaConsumer.ConsumeSingle(ctx)
		sCtx.Require().NotNil(msg, "Ожидается получение события о переходе заказа в статус canceled")
		sCtx.WithNewAttachment("canceled order message body", allure.JSON, msg.Value)

		err := json.Unmarshal(msg.Value, &orderEvent)
		sCtx.Require().NoError(err, "Ожидается успешная десериализация события")
		assert.OrderEvent(sCtx, domain.OrderEvent{OrderID: orderID, Status: domain.OrderStatusCancelled}, orderEvent)
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
