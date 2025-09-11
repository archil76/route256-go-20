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

func Test_OrderInfo(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	lomsServiseMock := mock.NewLomsServiseMock(ctrl)

	handler := NewServer(lomsServiseMock)

	orderCreateRequest := desc.OrderCreateRequest{
		UserId: tp.userID,
		Items: []*desc.Items{
			{
				Sku:   tp.sku,
				Count: tp.count,
			},
		},
	}
	orderID := int64(1)

	orderInfoResponseExpected := desc.OrderInfoResponse{
		Status:  string(model.AWAITINGPAYMENT),
		OrderId: orderID,
		Items: []*desc.Items{
			{
				Sku:   tp.sku,
				Count: tp.count,
			},
		},
	}

	orderExpected := model.Order{
		OrderID: orderID,
		UserID:  tp.userID,
		Status:  model.AWAITINGPAYMENT,
		Items: []model.Item{
			{
				Sku:   tp.sku,
				Count: tp.count,
			},
		},
	}

	orderInfoRequest := desc.OrderInfoRequest{
		OrderId: orderID,
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
		require.Equal(t, orderID, orderCreateResponse.OrderId)
	})

	t.Run("Информация по Заказу. Успешный путь", func(t *testing.T) {
		lomsServiseMock.OrderInfoMock.When(ctx, orderID).Then(&orderExpected, nil)

		orderInfoResponse, err := handler.OrderInfo(ctx, &orderInfoRequest)

		require.NoError(t, err)
		require.Equal(t, &orderInfoResponseExpected, orderInfoResponse)
	})

}
