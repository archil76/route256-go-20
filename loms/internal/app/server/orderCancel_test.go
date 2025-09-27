package server

import (
	"context"
	desc "route256/loms/internal/api"
	"route256/loms/internal/app/server/mock"
	"route256/loms/internal/domain/model"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func Test_OrderCancel(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	lomsServiseMock := mock.NewLomsServiseMock(ctrl)

	handler := NewServer(lomsServiseMock)

	orderCreateRequest := desc.OrderCreateRequest{
		UserID: tp.userID,
		Items: []*desc.Items{
			{
				Sku:   tp.sku,
				Count: tp.count,
			},
		},
	}
	orderID := int64(1)

	orderCancelRequest := desc.OrderCancelRequest{
		OrderID: orderID,
	}

	items := []model.Item{
		{
			Sku:   tp.sku,
			Count: tp.count,
		},
	}

	t.Run("Добавление Заказа. Успешный путь", func(t *testing.T) {
		lomsServiseMock.OrderCreateMock.When(ctx, tp.userID, items).Then(orderID, nil)

		orderCreateResponse, err := handler.OrderCreate(ctx, &orderCreateRequest)

		require.NoError(t, err)
		require.Equal(t, orderID, orderCreateResponse.OrderID)
	})

	t.Run("Отмена Заказа. Успешный путь", func(t *testing.T) {
		lomsServiseMock.OrderCancelMock.When(ctx, orderID).Then(nil)

		orderCancelResponse, err := handler.OrderCancel(ctx, &orderCancelRequest)

		require.NoError(t, err)
		require.Equal(t, &desc.OrderCancelResponse{}, orderCancelResponse)
	})
}
