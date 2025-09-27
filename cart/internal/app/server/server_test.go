package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"route256/cart/internal/domain/model"
	"strconv"

	"testing"
)

var (
	sku        = model.Sku(1076963)
	sku2       = model.Sku(1148162) // Должен быть больше sku для проверки сортировки получаемой корзины
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

	reportCart = model.ReportCart{
		Items: []model.ItemInСart{
			{SKU: sku, Count: count, Name: name, Price: price},
			{SKU: sku2, Count: count2, Name: name2, Price: price2},
		},
		TotalPrice: totalPrice,
	}
)

func TestHandler_All(t *testing.T) {
	t.Run("Test_ClearCartHandler", Test_ClearCartHandler)
	t.Run("Test_DeleteItemHandler", Test_DeleteItemHandler)
	t.Run("Test_AddItemHandler", Test_AddItemHandler)
	t.Run("Test_GetCartHandler", Test_GetCartHandler)
	t.Run("Test_CheckoutHandler", Test_CheckoutHandler)
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

func getCheckoutRequest(userID int64) (*http.Request, error) {
	request, err := http.NewRequest(
		http.MethodPost,
		"/checkout/{user_id]",
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
