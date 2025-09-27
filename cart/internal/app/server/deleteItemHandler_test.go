package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/domain/model"
	repo "route256/cart/internal/domain/repository/inmemoryrepository"
	cartsService "route256/cart/internal/domain/service"
	mock2 "route256/cart/internal/domain/service/mock"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func Test_DeleteItemHandler(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	cartRepositoryMock := mock2.NewCartsRepositoryMock(ctrl)
	productServiceMock := mock2.NewProductServiceMock(ctrl)
	lomsServiceMock := mock2.NewLomsServiceMock(ctrl)

	cartService := cartsService.NewCartsService(cartRepositoryMock, productServiceMock, lomsServiceMock)

	handler := NewServer(cartService)

	t.Run("Удаление sku из корзины. Успешный путь", func(t *testing.T) {
		t.Helper()

		cartRepositoryMock.DeleteItemMock.Expect(ctx, userID, model.Item{Sku: sku, Count: 0}).Return(&item, nil)

		request, err := getDeleteItemRequest(userID, sku)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.DeleteItem(writer, request)
		require.ErrorIs(t, nil, err)

		response := writer.Result()
		require.Equal(t, http.StatusNoContent, response.StatusCode)
	})

	t.Run("Удаление sku из корзины. Корзины нет", func(t *testing.T) {
		t.Helper()

		cartRepositoryMock.DeleteItemMock.Expect(ctx, userID, model.Item{Sku: sku, Count: 0}).Return(nil, repo.ErrCartDoesntExist)

		request, err := getDeleteItemRequest(userID, sku)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.DeleteItem(writer, request)
		require.ErrorIs(t, nil, err)

		response := writer.Result()
		require.Equal(t, http.StatusNoContent, response.StatusCode)
	})

	t.Run("Удаление sku из корзины. Невалидный sku", func(t *testing.T) {
		t.Helper()

		request, err := getDeleteItemRequest(userID, wrongSku)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.DeleteItem(writer, request)
		require.ErrorIs(t, nil, err)

		response := writer.Result()
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Удаление sku из корзины. Невалидный user", func(t *testing.T) {
		t.Helper()

		request, err := getDeleteItemRequest(wrongUserID, sku)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		handler.DeleteItem(writer, request)
		require.ErrorIs(t, nil, err)

		response := writer.Result()
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}
