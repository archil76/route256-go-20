package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/domain/model"
	repo "route256/cart/internal/domain/repository/inmemoryrepository"
	cartsService "route256/cart/internal/domain/service"
	mock2 "route256/cart/internal/domain/service/mock"
	"sort"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func Test_GetCartHandler(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	cartRepositoryMock := mock2.NewCartsRepositoryMock(ctrl)
	productServiceMock := mock2.NewProductServiceMock(ctrl)
	lomsServiceMock := mock2.NewLomsServiceMock(ctrl)

	cartService := cartsService.NewCartsService(cartRepositoryMock, productServiceMock, lomsServiceMock)
	skus := []model.Sku{sku, sku2}
	sort.Slice(skus, func(i, j int) bool { return skus[i] < skus[j] })

	products := []model.Product{*product, *product2}

	handler := NewServer(cartService)

	t.Run("Получение корзины. Успешный путь", func(t *testing.T) {
		t.Helper()

		productServiceMock.GetProductsBySkusMock.When(ctx, skus).Then(products, nil)
		productServiceMock.GetProductsBySkusMock.Optional()

		cartRepositoryMock.GetCartMock.Expect(ctx, userID).Return(cart, nil)

		request, err := getGetCartRequest(userID)
		require.NoError(t, err)

		writer := httptest.NewRecorder()

		handler.GetCart(writer, request)
		require.NoError(t, err)

		response := writer.Result()
		require.Equal(t, http.StatusOK, response.StatusCode)

		reportCart, err := decodeResponseBody(response)
		require.NoError(t, err)

		count32 := int32(count)           //nolint:gosec
		price32 := int32(price)           //nolint:gosec
		count232 := int32(count2)         //nolint:gosec
		price232 := int32(price2)         //nolint:gosec
		totalPrice32 := int32(totalPrice) //nolint:gosec

		wantedReportCart := ReportCart{
			Items: []ItemInСart{
				{
					SKU:   sku,
					Count: count32,
					Name:  name,
					Price: price32,
				},
				{
					SKU:   sku2,
					Count: count232,
					Name:  name2,
					Price: price232,
				},
			},
			TotalPrice: totalPrice32,
		}
		sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })
		require.Equal(t, wantedReportCart, reportCart)
	})

	t.Run("Получение корзины. Не валидный user", func(t *testing.T) {
		request, err := getGetCartRequest(wrongUserID)
		require.NoError(t, err)

		writer := httptest.NewRecorder()

		handler.GetCart(writer, request)
		require.NoError(t, err)

		response := writer.Result()
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Получение корзины. Корзины нет", func(t *testing.T) {
		cartRepositoryMock.GetCartMock.Expect(ctx, userID).Return(nil, repo.ErrCartDoesntExist)

		request, err := getGetCartRequest(userID)
		require.NoError(t, err)

		writer := httptest.NewRecorder()

		handler.GetCart(writer, request)
		require.NoError(t, err)

		response := writer.Result()
		require.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}
