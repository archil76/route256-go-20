package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	repo "route256/cart/internal/domain/repository/inmemoryrepository"
	cartsService "route256/cart/internal/domain/service"
	mock2 "route256/cart/internal/domain/service/mock"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func Test_ClearCartHandler(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	cartRepositoryMock := mock2.NewCartsRepositoryMock(ctrl)
	productServiceMock := mock2.NewProductServiceMock(ctrl)
	lomsServiceMock := mock2.NewLomsServiceMock(ctrl)

	cartService := cartsService.NewCartsService(cartRepositoryMock, productServiceMock, lomsServiceMock)

	handler := NewServer(cartService)

	t.Run("Очистка корзины. Успешный путь", func(t *testing.T) {
		request, err := getClearCartRequest(userID)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		cartRepositoryMock.DeleteItemsMock.Expect(ctx, userID).Return(userID, nil)

		handler.ClearCart(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusNoContent, response.StatusCode)
	})

	t.Run("Очистка корзины. Не валидный пользователь", func(t *testing.T) {
		request, err := getClearCartRequest(wrongUserID)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		cartRepositoryMock.DeleteItemsMock.Expect(ctx, userID).Return(userID, nil)

		handler.ClearCart(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Очистка корзины. Корзины нет", func(t *testing.T) {
		request, err := getClearCartRequest(userID)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		cartRepositoryMock.DeleteItemsMock.Expect(ctx, userID).Return(userID, repo.ErrCartDoesntExist)

		handler.ClearCart(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusNoContent, response.StatusCode)
	})

	t.Run("Очистка корзины. Корзина пуста", func(t *testing.T) {
		request, err := getClearCartRequest(userID)
		require.ErrorIs(t, nil, err)

		writer := httptest.NewRecorder()

		cartRepositoryMock.DeleteItemsMock.Expect(ctx, userID).Return(userID, nil)

		handler.ClearCart(writer, request)

		response := writer.Result()

		require.Equal(t, http.StatusNoContent, response.StatusCode)
	})

}
