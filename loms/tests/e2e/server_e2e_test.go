package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"route256/loms/internal/infra/config"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/asserts_wrapper/require"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ServerE struct {
	suite.Suite
	Host   string
	client *http.Client
}

func TestServerE(t *testing.T) {
	t.Parallel()
	suite.RunSuite(t, new(ServerE))
}

func (s *ServerE) BeforeAll(t provider.T) {
	env_var := os.Getenv("CONFIG_FILE")
	if env_var == "" {
		t.Fatalf("Не задана переменная окружения CONFIG_FILE")
		return
	}

	c, err := config.LoadConfig(os.Getenv("CONFIG_FILE"))
	if err != nil {
		t.Fatalf("Неверный формат конфига по адресу: %s", env_var)
		return
	}

	s.Host = fmt.Sprintf("http://%s:%s", c.Server.Host, c.Server.HttpPort)

	s.client = &http.Client{}

	t.Title("e2e test")
}

type CreateOrderRequest struct {
	UserID int64         `json:"userID"`
	Items  []ItemRequest `json:"items"`
}

type ItemRequest struct {
	Sku   int64  `json:"sku"`
	Count uint32 `json:"count"`
}

type OrderIDResponse struct {
	OrderID int64 `json:"orderID"`
}

func (s *ServerE) TestServerE(t provider.T) {
	t.Parallel()
	t.Helper()

	sku := int64(1076963)
	sku2 := int64(1148162) // должен быть больше Sku для проверки сортировки получаемой корзины
	wrongSku := int64(1076963000)

	count := uint32(2)
	count2 := uint32(3)
	userID := int64(1022222)

	createOrderRequest := CreateOrderRequest{
		UserID: userID,
		Items: []ItemRequest{
			{
				Sku:   sku,
				Count: count,
			},
			{
				Sku:   sku2,
				Count: count2,
			},
		}}

	createWrongOrderRequest := CreateOrderRequest{
		UserID: userID,
		Items: []ItemRequest{
			{
				Sku:   sku,
				Count: count,
			},
			{
				Sku:   wrongSku,
				Count: count2,
			},
		}}

	t.Title("Создание заказа")

	t.WithNewStep("Действие: Создание заказа", func(t provider.StepCtx) {

		request, err := getCreateOrderRequest(s.Host, createOrderRequest)
		require.NoError(t, err)

		response, err := s.client.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		orderIDResponse, err := decodeResponseBody(response, OrderIDResponse{})
		require.NoError(t, err)
		if value, ok := orderIDResponse.(OrderIDResponse); ok {
			require.Greater(t, 0, value.OrderID)
		} else {
			require.False(t, ok)
		}

	})

	t.WithNewStep("Действие: Не валидный sku", func(t provider.StepCtx) {

		request, err := getCreateOrderRequest(s.Host, createWrongOrderRequest)
		require.NoError(t, err)

		response, err := s.client.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)

	})
}

func getCreateOrderRequest(host string, createOrderRequest CreateOrderRequest) (*http.Request, error) {
	body, err := json.Marshal(createOrderRequest)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(body)
	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/order/create", host),
		reader,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

func decodeResponseBody(response *http.Response, structResponse interface{}) (interface{}, error) {
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	err := decoder.Decode(&structResponse)

	return structResponse, err
}

func StatusCode(sCtx provider.StepCtx, expected, actual int) {
	sCtx.Require().Equal(expected, actual, "Не совпадает статус кодР")
}
