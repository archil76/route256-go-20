package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	producrepo "route256/cart/internal/domain/repository/productservicerepository"
	cartsService "route256/cart/internal/domain/service"
	mock2 "route256/cart/internal/domain/service/mock"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func Test_AddItemHandler(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	cartRepositoryMock := mock2.NewCartsRepositoryMock(ctrl)
	productServiceMock := mock2.NewProductServiceMock(ctrl)
	lomsServiceMock := mock2.NewLomsServiceMock(ctrl)

	cartService := cartsService.NewCartsService(cartRepositoryMock, productServiceMock, lomsServiceMock)

	handler := NewServer(cartService)

	t.Run("Добавление sku в корзину. Успешный путь", func(t *testing.T) {
		productServiceMock.GetProductBySkuMock.When(ctx, sku).Then(product, nil)
		productServiceMock.GetProductBySkuMock.When(ctx, sku2).Then(product2, nil)

		cartRepositoryMock.AddItemMock.When(ctx, userID, item).Then(&item, nil)
		cartRepositoryMock.AddItemMock.When(ctx, userID, item2).Then(&item2, nil)

		lomsServiceMock.StockInfoMock.When(ctx, sku).Then(265, nil)
		lomsServiceMock.StockInfoMock.When(ctx, sku2).Then(300, nil)

		request, err := getAddItemRequest(addItemRequest, userID, sku)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.AddItem(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusOK, response.StatusCode)

		request2, err := getAddItemRequest(addItemRequest2, userID, sku2)
		require.NoError(t, err)

		handler.AddItem(writer, request2)
		require.NoError(t, err)

		response = writer.Result()

		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("Добавление sku в корзину. Не валидный user", func(t *testing.T) {

		request, err := getAddItemRequest(addItemRequest, wrongUserID, sku)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.AddItem(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)

	})

	t.Run("Добавление sku в корзину. Не валидный sku", func(t *testing.T) {

		request, err := getAddItemRequest(addItemRequest, userID, wrongSku)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.AddItem(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)

	})

	t.Run("Добавление sku в корзину. Не существующий product", func(t *testing.T) {
		productServiceMock.GetProductBySkuMock.When(ctx, 100).Then(nil, producrepo.ErrProductNotFound)

		request, err := getAddItemRequest(addItemRequest, userID, 100)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.AddItem(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusPreconditionFailed, response.StatusCode)
	})

	t.Run("Добавление sku в корзину. Не хватает остатка", func(t *testing.T) {
		productServiceMock.GetProductBySkuMock.When(ctx, 139275865).Then(product, nil)

		lomsServiceMock.StockInfoMock.When(ctx, 139275865).Then(0, nil)

		request, err := getAddItemRequest(addItemRequest, userID, 139275865)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.AddItem(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusPreconditionFailed, response.StatusCode)
	})
}
