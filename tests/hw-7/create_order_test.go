package hw7

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
	"route256/tests/app/domain"
)

func (s *Suite) TestCreateOrderSuccess(t provider.T) {
	t.Title("Успешное создание заказа")

	const (
		userID = int64(1000)
		sku    = int32(1076963)
		count  = int64(5)
	)

	var (
		orderID int64
		ctx     = context.Background()
	)

	t.WithNewParameters(
		"userID", userID,
		"sku", sku,
		"count", count,
	)

	t.WithNewStep("Создание заказа через апи", func(sCtx provider.StepCtx) {
		// создание заказа через метод апи
		var statusCode int
		orderID, statusCode = s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: count},
		})

		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderID(sCtx, orderID)

		var orderEvent domain.OrderEvent

		// валидация события об успешном создании заказа
		msg := s.kafkaConsumer.ConsumeSingle(ctx)
		sCtx.Require().NotNil(msg, "Ожидается получение события о создании заказа")
		sCtx.WithNewAttachment("new order message body", allure.JSON, msg.Value)

		err := json.Unmarshal(msg.Value, &orderEvent)
		sCtx.Require().NoError(err, "Ожидается успешная десериализация события")
		assert.OrderEvent(sCtx, domain.OrderEvent{OrderID: orderID, Status: domain.OrderStatusNew}, orderEvent)

		// валидация события об успешном бронировании и переходе в статус ожидания оплаты
		msg = s.kafkaConsumer.ConsumeSingle(ctx)
		sCtx.Require().NotNil(msg, "Ожидается получение события о переходе заказа в статус awaiting payment")
		sCtx.WithNewAttachment("awaiting payment message body", allure.JSON, msg.Value)

		err = json.Unmarshal(msg.Value, &orderEvent)
		sCtx.Require().NoError(err, "Ожидается успешная десериализация события")
		assert.OrderEvent(sCtx, domain.OrderEvent{OrderID: orderID, Status: domain.OrderStatusAwaitingPayment}, orderEvent)
	})

	t.WithNewStep("Проверка статуса заказа", func(sCtx provider.StepCtx) {
		res, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, orderID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.OrderStatus(sCtx, domain.OrderStatusAwaitingPayment, res.Status)
	})
}

func (s *Suite) TestCreateFailed(t provider.T) {
	t.Title("Неуспешная попытка резервирования стоков при созданиии заказа")

	const (
		userID = int64(1000)
		sku    = int32(1076963)
		count  = int64(5000000)
	)

	var (
		ctx = context.Background()
	)

	t.WithNewParameters(
		"userID", userID,
		"sku", sku,
		"count", count,
	)

	t.WithNewStep("Создание заказа через апи", func(sCtx provider.StepCtx) {
		// создание заказа через метод апи
		_, statusCode := s.lomsClient.OrderCreate(ctx, sCtx, userID, []domain.OrderItem{
			{Sku: sku, Count: count},
		})

		// {"code":9, "message":"invalid stock", "details":[]}
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)

		// валидация события об успешном создании заказа
		msg := s.kafkaConsumer.ConsumeSingle(ctx)
		sCtx.Require().NotNil(msg, "Ожидается получение события о создании заказа")
		sCtx.WithNewAttachment("new order message body", allure.JSON, msg.Value)

		var orderEvent domain.OrderEvent

		err := json.Unmarshal(msg.Value, &orderEvent)
		sCtx.Require().NoError(err, "Ожидается успешная десериализация события")
		sCtx.Require().True(orderEvent.OrderID > 0, "Ожидается, что OrderID в событии будет больше 0")
		sCtx.Require().Equal(domain.OrderStatusNew, orderEvent.Status, "Ожидается статус нового заказа")

		orderId := orderEvent.OrderID

		// валидация события о переходе заказа в статус failed
		msg = s.kafkaConsumer.ConsumeSingle(ctx)
		sCtx.Require().NotNil(msg, "Ожидается получение события о переходе заказа в статус failed")
		sCtx.WithNewAttachment("failed order message body", allure.JSON, msg.Value)

		err = json.Unmarshal(msg.Value, &orderEvent)
		sCtx.Require().NoError(err, "Ожидается успешная десериализация события")
		assert.OrderEvent(sCtx, domain.OrderEvent{OrderID: orderId, Status: domain.OrderStatusFailed}, orderEvent)
	})
}
