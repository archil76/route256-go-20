package hw7

import (
	"context"
	"encoding/json"
	"net/http"
	"route256/tests/app/assert"
	"route256/tests/app/domain"

	"github.com/ozontech/allure-go/pkg/allure"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestOrderPaySuccess(t provider.T) {
	t.Title("Успешная оплата заказа")

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

	t.WithNewStep("Проверка что заказ уже перешел в awaiting payment", func(sCtx provider.StepCtx) {
		res, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderStatus(sCtx, domain.OrderStatusAwaitingPayment, res.Status)
	})

	t.WithNewStep("Оплата заказа", func(sCtx provider.StepCtx) {
		statusCode := s.lomsClient.OrderPay(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		msg := s.kafkaConsumer.ConsumeSingle(ctx)
		sCtx.Require().NotNil(msg, "Ожидается получения сообщения о переходе заказа в статус paid")
		sCtx.WithNewAttachment("paid message body", allure.JSON, msg.Value)

		err := json.Unmarshal(msg.Value, &orderEvent)
		sCtx.Require().NoError(err, "Ожидается успешная десериализация события")
		assert.OrderEvent(sCtx, domain.OrderEvent{OrderID: orderID, Status: domain.OrderStatusPaid}, orderEvent)
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
