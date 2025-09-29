package server

import (
	"context"
	desc "route256/loms/internal/api"
	"route256/loms/internal/app/server/mock"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func Test_StocksInfo(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	lomsServiseMock := mock.NewLomsServiseMock(ctrl)

	handler := NewServer(lomsServiseMock)

	stocksInfoRequest := desc.StocksInfoRequest{
		Sku: tp.sku,
	}
	stocksInfoResponseExpected := desc.StocksInfoResponse{
		Count: tp.countInStock,
	}

	t.Run("Информация о стоке. Успешный путь", func(t *testing.T) {
		lomsServiseMock.StocksInfoMock.When(ctx, tp.sku).Then(tp.countInStock, nil)

		stocksInfoResponse, err := handler.StocksInfo(ctx, &stocksInfoRequest)

		require.NoError(t, err)
		require.Equal(t, &stocksInfoResponseExpected, stocksInfoResponse)
	})

}
