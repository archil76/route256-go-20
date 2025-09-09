package service

import (
	"context"
	"route256/loms/internal/domain/model"
	orderrepo "route256/loms/internal/domain/repository/inmemoryrepository/order"
	stockrepo "route256/loms/internal/domain/repository/inmemoryrepository/stock"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_OrderCreate(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()

	orderRepository := orderrepo.NewOrderInMemoryRepository(10, &counter)

	stockRepository := stockrepo.NewStockInMemoryRepository(10, source)

	handler := NewLomsService(orderRepository, stockRepository)

	items := []model.Item{
		{
			Sku:   tp.sku,
			Count: tp.count,
		},
	}

	t.Run("Добавление Заказа. Успешный путь", func(t *testing.T) {

		orderID, err := handler.OrderCreate(ctx, 10000, items)

		require.NoError(t, err)
		require.NotEqual(t, 0, orderID)

		orderID2, err := handler.OrderCreate(ctx, 10000, items)

		require.NoError(t, err)
		require.Greater(t, orderID2, orderID)

	})

}
