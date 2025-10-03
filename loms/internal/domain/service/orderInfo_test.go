package service

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

	t.Run("Получение Заказа. Проверка", func(t *testing.T) {
		orderExpected1 := model.Order{
			OrderID: orderID,
			UserID:  10000,
			Status:  model.AWAITINGPAYMENT,
			Items: []model.Item{
				{Sku: tp.sku, Count: tp.count},
			},
		}

		orderExpected2 := model.Order{
			OrderID: orderID2,
			UserID:  10000,
			Status:  model.AWAITINGPAYMENT,
			Items: []model.Item{
				{Sku: tp.sku, Count: tp.count},
			},
		}

		order1, err := handler.OrderInfo(ctx, orderID)
		require.NoError(t, err)

		require.Equal(t, orderExpected1, *order1)

		order2, err := handler.OrderInfo(ctx, orderID2)
		require.NoError(t, err)
		require.Equal(t, orderExpected2, *order2)

	})

}
