package loms

import (
	"context"
	"route256/loms/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_OrderInfo(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()

	testHandler := NewLomsServiceWithMock(t)

	t.Run("Проверка статуса. Успешный путь", func(t *testing.T) {
		handler := testHandler.handler
		testHandler.orderRepositoryMock.GetByIDMock.When(ctx, tp.orderID).Then(&canceledOrder, nil)

		orderRes, err := handler.OrderInfo(ctx, tp.orderID)
		require.NoError(t, err)
		require.Equal(t, model.CANCELED, orderRes.Status)

		testHandler.orderRepositoryMock.GetByIDMock.When(ctx, tp.orderID2).Then(&awaitingPaymentOrder, nil)

		order2Res, err := handler.OrderInfo(ctx, tp.orderID2)
		require.NoError(t, err)
		require.Equal(t, model.AWAITINGPAYMENT, order2Res.Status)
	})
}
