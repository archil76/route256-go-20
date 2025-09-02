package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/domain/model"
	"sort"
	"strconv"

	repo "route256/cart/internal/domain/repository/inmemoryrepository"
	producrepo "route256/cart/internal/domain/repository/productservicerepository"
	cartsService "route256/cart/internal/domain/service"
	mock2 "route256/cart/internal/domain/service/mock"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

var (
	sku        = model.Sku(1076963)
	sku2       = model.Sku(1148162) // должен быть больше sku для проверки сортировки получаемой корзины
	name       = "Flashlight"
	name2      = "Greenhouse"
	count      = uint32(2)
	count2     = uint32(3)
	userID     = int64(2009999999)
	price      = uint32(1000)
	price2     = uint32(2000)
	totalPrice = price*count + price2*count2

	wrongUserID = int64(-111111)
	wrongSku    = model.Sku(-1076963)

	product  = &model.Product{Sku: sku, Name: name, Price: price}
	product2 = &model.Product{Sku: sku2, Name: name2, Price: price2}
	cart     = &model.Cart{UserID: userID, Items: map[model.Sku]uint32{sku: count, sku2: count2}}
	//cart2    = &model.Cart{UserID: userID, Items: map[model.Sku]uint32{sku2: count2}}
	item  = model.Item{Sku: sku, Count: count}
	item2 = model.Item{Sku: sku2, Count: count2}

	addItemRequest = AddItemRequest{
		Count: int32(count),
	}

	addItemRequest2 = AddItemRequest{
		Count: int32(count2),
	}
)

func TestHandler_All(t *testing.T) {
	t.Run("Test_ClearCartHandler", Test_ClearCartHandler)
	t.Run("Test_DeleteItemHandler", Test_DeleteItemHandler)
	t.Run("Test_AddItemHandler", Test_AddItemHandler)
	t.Run("Test_GetCartHandler", Test_GetCartHandler)
}

func Test_ClearCartHandler(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	cartRepositoryMock := mock2.NewCartsRepositoryMock(ctrl)
	productServiceMock := mock2.NewProductServiceMock(ctrl)

	cartService := cartsService.NewCartsService(cartRepositoryMock, productServiceMock)

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

func Test_DeleteItemHandler(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	cartRepositoryMock := mock2.NewCartsRepositoryMock(ctrl)
	productServiceMock := mock2.NewProductServiceMock(ctrl)

	cartService := cartsService.NewCartsService(cartRepositoryMock, productServiceMock)

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

func Test_AddItemHandler(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	cartRepositoryMock := mock2.NewCartsRepositoryMock(ctrl)
	productServiceMock := mock2.NewProductServiceMock(ctrl)

	cartService := cartsService.NewCartsService(cartRepositoryMock, productServiceMock)

	handler := NewServer(cartService)

	t.Run("Добавление sku в корзину. Успешный путь", func(t *testing.T) {
		productServiceMock.GetProductBySkuMock.When(ctx, sku).Then(product, nil)
		productServiceMock.GetProductBySkuMock.When(ctx, sku2).Then(product2, nil)

		cartRepositoryMock.AddItemMock.When(ctx, userID, item).Then(&item, nil)
		cartRepositoryMock.AddItemMock.When(ctx, userID, item2).Then(&item2, nil)

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
}

func Test_GetCartHandler(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()
	ctrl := minimock.NewController(t)

	cartRepositoryMock := mock2.NewCartsRepositoryMock(ctrl)
	productServiceMock := mock2.NewProductServiceMock(ctrl)

	cartService := cartsService.NewCartsService(cartRepositoryMock, productServiceMock)

	handler := NewServer(cartService)

	t.Run("Получение корзины. Успешный путь", func(t *testing.T) {
		t.Helper()

		productServiceMock.GetProductBySkuMock.When(ctx, sku).Then(product, nil)
		productServiceMock.GetProductBySkuMock.When(ctx, sku2).Then(product2, nil)
		productServiceMock.GetProductBySkuMock.Optional()

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

func getDeleteItemRequest(userID int64, sku model.Sku) (*http.Request, error) {
	request, err := http.NewRequest(
		http.MethodDelete,
		"/user/{user_id]/cart/{sku_id}",
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.SetPathValue("user_id", strconv.FormatInt(userID, 10))
	request.SetPathValue("sku_id", strconv.FormatInt(sku, 10))

	return request, nil
}

func getAddItemRequest(addItemRequest AddItemRequest, userID int64, sku model.Sku) (*http.Request, error) {
	body, err := json.Marshal(addItemRequest)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(body)
	request, err := http.NewRequest(
		http.MethodPost,
		"/user/{user_id]/cart/{sku_id}",
		reader,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.SetPathValue("user_id", strconv.FormatInt(userID, 10))
	request.SetPathValue("sku_id", strconv.FormatInt(sku, 10))

	return request, nil
}

func getGetCartRequest(userID int64) (*http.Request, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		"/user/{user_id]/cart",
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.SetPathValue("user_id", strconv.FormatInt(userID, 10))

	return request, nil
}

func getClearCartRequest(userID int64) (*http.Request, error) {
	request, err := http.NewRequest(
		http.MethodDelete,
		"/user/{user_id]/cart",
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.SetPathValue("user_id", strconv.FormatInt(userID, 10))

	return request, nil
}

func decodeResponseBody(response *http.Response) (ReportCart, error) {
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	reportCart := ReportCart{}

	err := decoder.Decode(&reportCart)

	return reportCart, err
}
