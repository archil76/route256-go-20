package loms

import (
	"context"
	"route256/loms/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_OrderCreate(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()

	testHandler := NewLomsServiceWithMock(t)

	t.Run("Добавление Заказа. Успешный путь", func(t *testing.T) {
		var err error
		var orderIDRes, orderID2Res int64

		testHandler.orderRepositoryMock.CreateMock.When(ctx, newOrder).Then(&order, nil)
		testHandler.outboxServiceMock.CreateMessageMock.When(ctx, order.OrderID, model.NEWSTATUS)
		testHandler.stockRepositoryMock.ReserveMock.When(ctx, items).Then(stocks, nil)
		testHandler.orderRepositoryMock.SetStatusMock.When(ctx, order, model.AWAITINGPAYMENT).Then(nil)
		testHandler.outboxServiceMock.CreateMessageMock.When(ctx, order.OrderID, model.AWAITINGPAYMENT)

		handler := testHandler.handler

		orderIDRes, err = handler.OrderCreate(ctx, tp.userID, items)

		require.NoError(t, err)
		require.Equal(t, tp.orderID, orderIDRes)

		testHandler.orderRepositoryMock.CreateMock.When(ctx, newOrder2).Then(&order2, nil)
		testHandler.outboxServiceMock.CreateMessageMock.When(ctx, order2.OrderID, model.NEWSTATUS)
		testHandler.orderRepositoryMock.SetStatusMock.When(ctx, order2, model.AWAITINGPAYMENT).Then(nil)
		testHandler.outboxServiceMock.CreateMessageMock.When(ctx, order2.OrderID, model.AWAITINGPAYMENT)

		orderID2Res, err = handler.OrderCreate(ctx, tp.userID2, items)

		require.NoError(t, err)
		require.Equal(t, tp.orderID2, orderID2Res)
	})
}
