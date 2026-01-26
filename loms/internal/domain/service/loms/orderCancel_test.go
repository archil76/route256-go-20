package loms

import (
	"context"
	"route256/loms/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_OrderCancel(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()

	testHandler := NewLomsServiceWithMock(t)

	t.Run("Отмена Заказа. Успешный путь", func(t *testing.T) {

		testHandler.orderRepositoryMock.GetByIDMock.When(ctx, tp.orderID).Then(&awaitingPaymentOrder, nil)

		testHandler.stockRepositoryMock.ReserveCancelMock.When(ctx, items).Then(nil)
		testHandler.orderRepositoryMock.SetStatusMock.When(ctx, awaitingPaymentOrder, model.CANCELED).Then(nil)
		testHandler.outboxServiceMock.CreateMessageMock.When(ctx, tp.orderID, model.CANCELED)

		handler := testHandler.handler

		err := handler.OrderCancel(ctx, tp.orderID)

		require.NoError(t, err)
	})
}
