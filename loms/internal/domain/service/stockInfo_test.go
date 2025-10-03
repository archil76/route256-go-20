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

	handler := NewLomsServiceWithInMemoryRepository()

	t.Run("Информация о стоке. Успешный путь", func(t *testing.T) {

		count, err := handler.StocksInfo(ctx, 139275865)

		require.NoError(t, err)
		require.NotEqual(t, 0, count)

	})

}
