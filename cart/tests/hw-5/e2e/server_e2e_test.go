package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"route256/cart/internal/infra/config"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/asserts_wrapper/require"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ServerE struct {
	suite.Suite
	Host   string
	client *http.Client
	ctx    context.Context
}

func TestServerE(t *testing.T) {
	t.Parallel()
	suite.RunSuite(t, new(ServerE))
}

func (s *ServerE) BeforeAll(t provider.T) {
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

func (s *ServerE) TestServerParallel(t provider.T) {

	t.Helper()

	userID := int64(1022221)
	skus := []int64{
		1076963,
		1148162,
		1625903,
		2618151,
		2956315,
		2958025,
		3596599,
		4465995,
		4288068,
	}
	t.Title("Проверка получения корзины")

	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Подготовка: Очистка корзины", func(t provider.StepCtx) {

			request, err := getClearCartRequest(ctx, s.Host, userID)
			require.NoError(t, err)

			response, err := s.client.Do(request)
			require.NoError(t, err)
			require.Equal(t, http.StatusNoContent, response.StatusCode)
		})

		t.WithNewStep("Подготовка: Проверка что корзина пуста", func(t provider.StepCtx) {
			request, err := getGetCartRequest(ctx, s.Host, userID)
			require.NoError(t, err)

			response, err := s.client.Do(request)
			require.NoError(t, err)

			require.Equal(t, http.StatusNotFound, response.StatusCode)
		})

		t.WithNewStep("Подготовка: Наполнение корзины", func(t provider.StepCtx) {

			for _, sku := range skus {
				request, err := getAddItemRequest(ctx, s.Host, testAddItemRequest{
					Count: 3,
				}, userID, sku)
				require.NoError(t, err)

				response, err := s.client.Do(request)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, response.StatusCode)

			}
		})
	})

	t.WithNewStep("Действие: Получение", func(t provider.StepCtx) {

		request, err := getGetCartRequest(ctx, s.Host, userID)
		require.NoError(t, err)

		response, err := s.client.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		reportCart, err := decodeResponseBody(response)
		require.NoError(t, err)

		sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })

		require.Equal(t, len(skus), len(reportCart.Items))
	})

}

func (s *ServerE) TestServerParallelWrongSku(t provider.T) {

	t.Helper()

	userID := int64(1022222)
	skus := []int64{
		32638658,
		32605854,
		32205848,
	}
	t.Title("Проверка получения корзины")

	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Подготовка: Очистка корзины", func(t provider.StepCtx) {

			request, err := getClearCartRequest(ctx, s.Host, userID)
			require.NoError(t, err)

			response, err := s.client.Do(request)
			require.NoError(t, err)
			require.Equal(t, http.StatusNoContent, response.StatusCode)
		})

		t.WithNewStep("Подготовка: Проверка что корзина пуста", func(t provider.StepCtx) {
			request, err := getGetCartRequest(ctx, s.Host, userID)
			require.NoError(t, err)

			response, err := s.client.Do(request)
			require.NoError(t, err)

			require.Equal(t, http.StatusNotFound, response.StatusCode)
		})

		t.WithNewStep("Подготовка: Наполнение корзины", func(t provider.StepCtx) {
			for _, sku := range skus {
				request, err := getAddItemRequest(ctx, s.Host, testAddItemRequest{
					Count: 3,
				}, userID, sku)
				require.NoError(t, err)

				response, err := s.client.Do(request)

				require.NoError(t, err)
				require.Equal(t, http.StatusPreconditionFailed, response.StatusCode)
			}
		})
	})

	t.WithNewStep("Действие: Получение", func(t provider.StepCtx) {

		request, err := getGetCartRequest(ctx, s.Host, userID)
		require.NoError(t, err)

		response, err := s.client.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, response.StatusCode)

	})

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
