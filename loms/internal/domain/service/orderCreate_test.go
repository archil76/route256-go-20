package service

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

	handler := NewLomsServiceWithInMemoryRepository()

	items := []model.Item{
		{
			Sku:   tp.sku,
			Count: tp.count,
		},
	}
	var orderID, orderID2 int64
	t.Run("Добавление Заказа. Успешный путь", func(t *testing.T) {
		var err error
		orderID, err = handler.OrderCreate(ctx, 10000, items)

		require.NoError(t, err)
		require.NotEqual(t, 0, orderID)

		orderID2, err = handler.OrderCreate(ctx, 10000, items)

		require.NoError(t, err)
		require.Greater(t, orderID2, orderID)

		order1, err := handler.OrderInfo(ctx, orderID)
		require.NoError(t, err)
		require.Equal(t, model.AWAITINGPAYMENT, order1.Status)

		order2, err := handler.OrderInfo(ctx, orderID2)
		require.NoError(t, err)
		require.Equal(t, model.AWAITINGPAYMENT, order2.Status)

	})

	t.Run("Добавление Заказа. Проверка статуса", func(t *testing.T) {

		order1, err := handler.OrderInfo(ctx, orderID)
		require.NoError(t, err)
		require.Equal(t, model.AWAITINGPAYMENT, order1.Status)

		order2, err := handler.OrderInfo(ctx, orderID2)
		require.NoError(t, err)
		require.Equal(t, model.AWAITINGPAYMENT, order2.Status)

	})
}
