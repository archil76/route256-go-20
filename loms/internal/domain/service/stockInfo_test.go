package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StockInfo(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()

	testHandler := NewLomsServiceWithMock(t)

	t.Run("Проверка статуса. Успешный путь", func(t *testing.T) {
		handler := testHandler.handler
		var count uint32
		var err error

		testHandler.stockRepositoryMock.GetBySKUMock.When(ctx, tp.sku).Then(tp.count, nil)

		count, err = handler.StocksInfo(ctx, tp.sku)
		require.NoError(t, err)
		require.Equal(t, tp.count, count)
	})
}
