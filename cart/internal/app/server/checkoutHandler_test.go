package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/domain/model"
	cartsService "route256/cart/internal/domain/service"
	mock2 "route256/cart/internal/domain/service/mock"
	"sort"
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
	skusASC := []model.Sku{sku, sku2}
	sort.Slice(skusASC, func(i, j int) bool { return skusASC[i] < skusASC[j] })
	skusDESC := []model.Sku{sku, sku2}
	sort.Slice(skusDESC, func(i, j int) bool { return skusDESC[i] >= skusDESC[j] })

	products := []model.Product{*product, *product2}

	handler := NewServer(cartService)

	t.Run("Оформление заказа по корзине. Успешный путь", func(t *testing.T) {
		productServiceMock.GetProductsBySkusMock.When(ctx, skusASC).Then(products, nil)
		productServiceMock.GetProductsBySkusMock.When(ctx, skusDESC).Then(products, nil)
		productServiceMock.GetProductsBySkusMock.Optional()

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
