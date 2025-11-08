package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"route256/cart/internal/infra/config"
	"strconv"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type Suite struct {
	suite.Suite
	Host   string
	client *http.Client
	ctx    context.Context
}

func TestSuite(t *testing.T) {
	suite.RunNamedSuite(t, "Cart Домашнее задание 5", new(Suite))
}

func (s *Suite) BeforeAll(t provider.T) {
	s.ctx = context.Background()

	envVar := os.Getenv("CONFIG_FILE")
	if envVar == "" {
		t.Fatalf("Не задана переменная окружения CONFIG_FILE")
		return
	}

	c, err := config.LoadConfig(os.Getenv("CONFIG_FILE"))
	if err != nil {
		t.Fatalf("Неверный формат конфига по адресу: %s", envVar)
		return
	}

	s.Host = fmt.Sprintf("http://%s:%s", c.Server.Host, c.Server.Port)

	s.client = &http.Client{}

	t.Title("e2e test")
}

type testReportCart struct {
	Items      []testItemInСart
	TotalPrice uint32
}

type testItemInСart struct {
	SKU   int64
	Count uint32
	Name  string
	Price uint32
}

type testAddItemRequest struct {
	Count int32
}

func getAddItemRequest(ctx context.Context, host string, addItemRequest testAddItemRequest, userID int64, sku int64) (*http.Request, error) {
	body, err := json.Marshal(addItemRequest)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(body)
	request, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		fmt.Sprintf("%s/user/%s/cart/%s", host, strconv.FormatInt(userID, 10), strconv.FormatInt(sku, 10)),
		reader,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

func getGetCartRequest(ctx context.Context, host string, userID int64) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("%s/user/%s/cart", host, strconv.FormatInt(userID, 10)),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

func getClearCartRequest(ctx context.Context, host string, userID int64) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/user/%s/cart", host, strconv.FormatInt(userID, 10)),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

func decodeResponseBody(response *http.Response) (testReportCart, error) {
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	reportCart := testReportCart{}

	err := decoder.Decode(&reportCart)

	return reportCart, err
}
