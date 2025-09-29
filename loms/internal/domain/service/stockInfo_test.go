package service

import (
	"context"
	orderrepo "route256/loms/internal/domain/repository/inmemoryrepository/order"
	stockrepo "route256/loms/internal/domain/repository/inmemoryrepository/stock"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StockInfo(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()

	orderRepository := orderrepo.NewOrderInMemoryRepository(10, &counter)

	stockRepository := stockrepo.NewStockInMemoryRepository(10)

	handler := NewLomsService(orderRepository, stockRepository)

	t.Run("Информация о стоке. Успешный путь", func(t *testing.T) {

		count, err := handler.StocksInfo(ctx, 139275865)

		require.NoError(t, err)
		require.NotEqual(t, 0, count)

	})

}
