package loms

import (
	"context"
	"route256/loms/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_OrderPay(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()

	testHandler := NewLomsServiceWithMock(t)

	t.Run("Оплата Заказа. Успешный путь", func(t *testing.T) {
		testHandler.orderRepositoryMock.GetByIDMock.When(ctx, tp.orderID).Then(&awaitingPaymentOrder, nil)
		testHandler.stockRepositoryMock.ReserveRemoveMock.When(ctx, items).Then(nil)
		testHandler.orderRepositoryMock.SetStatusMock.When(ctx, awaitingPaymentOrder, model.PAYED).Then(nil)
		testHandler.outboxServiceMock.CreateMessageMock.When(ctx, awaitingPaymentOrder.OrderID, model.PAYED)
		testHandler.poolerMock.InTxMock.Set(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})
		handler := testHandler.handler

		err := handler.OrderPay(ctx, tp.orderID)

		require.NoError(t, err)
	})

	t.Run("Оплата отмененного заказа.", func(t *testing.T) {
		testHandler.orderRepositoryMock.GetByIDMock.When(ctx, tp.orderID2).Then(&canceledOrder, nil)

		handler := testHandler.handler

		err := handler.OrderPay(ctx, tp.orderID2)

		require.ErrorIs(t, model.ErrInvalidOrderStatus, err)
	})

}
