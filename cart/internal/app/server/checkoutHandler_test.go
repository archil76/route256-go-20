package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	cartsService "route256/cart/internal/domain/service"
	mock2 "route256/cart/internal/domain/service/mock"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func Test_CheckoutHandler(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	cartRepositoryMock := mock2.NewCartsRepositoryMock(ctrl)
	productServiceMock := mock2.NewProductServiceMock(ctrl)
	lomsServiceMock := mock2.NewLomsServiceMock(ctrl)

	cartService := cartsService.NewCartsService(cartRepositoryMock, productServiceMock, lomsServiceMock)

	handler := NewServer(cartService)

	t.Run("Оформление заказа по корзине. Успешный путь", func(t *testing.T) {
		productServiceMock.GetProductBySkuMock.When(ctx, sku).Then(product, nil)
		productServiceMock.GetProductBySkuMock.When(ctx, sku2).Then(product2, nil)
		productServiceMock.GetProductBySkuMock.Optional()

		cartRepositoryMock.GetCartMock.Expect(ctx, userID).Return(cart, nil)
		cartRepositoryMock.DeleteItemsMock.Expect(ctx, userID).Return(userID, nil)

		lomsServiceMock.OrderCreateMock.When(ctx, userID, &reportCart).Then(1, nil)

		request, err := getCheckoutRequest(userID)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.Checkout(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("Оформление заказа по корзине. Не валидный user", func(t *testing.T) {

		request, err := getCheckoutRequest(wrongUserID)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.Checkout(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)

	})
}
