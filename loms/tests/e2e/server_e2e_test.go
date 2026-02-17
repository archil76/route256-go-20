//go:build e2e_test

package hw3e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type Suite struct {
	suite.Suite

	cartClient *Client
	lomsClient *Client
}

func TestSuite(t *testing.T) {
	suite.RunNamedSuite(t, "hw 3 student", new(Suite))
}

func (s *Suite) BeforeAll(t provider.T) {

	s.cartClient = NewClient("http://localhost:8080")
	s.lomsClient = NewClient("http://localhost:8084")
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
