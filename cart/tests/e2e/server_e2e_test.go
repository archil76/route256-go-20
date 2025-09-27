//go:build e2e_test

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"route256/cart/internal/infra/config"
	"sort"
	"strconv"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
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

func (s *ServerE) TestServerE(t provider.T) {
	t.Parallel()
	t.Helper()

	sku := int64(1076963)
	sku2 := int64(1148162) // должен быть больше sku для проверки сортировки получаемой корзины
	wrongSku := uint32(1076963000)

	count := int32(2)
	count2 := int32(3)
	userID := int64(1022222)

	countUint32 := uint32(count)   //nolint:gosec
	count2Uint32 := uint32(count2) //nolint:gosec

	addItemRequest := testAddItemRequest{
		Count: count,
	}
	addItemRequest2 := testAddItemRequest{
		Count: count2,
	}

	t.Title("Проверка удаления товара из корзины")

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Подготовка: Очистка корзины", func(t provider.StepCtx) {

			request, err := getClearCartRequest(s.Host, userID)
			require.NoError(t, err)

			response, err := s.client.Do(request)
			require.NoError(t, err)
			require.Equal(t, http.StatusNoContent, response.StatusCode)
		})

		t.WithNewStep("Подготовка: Проверка что корзина пуста", func(t provider.StepCtx) {
			request, err := getGetCartRequest(s.Host, userID)
			require.NoError(t, err)

			response, err := s.client.Do(request)
			require.NoError(t, err)

			require.Equal(t, http.StatusNotFound, response.StatusCode)
		})

		t.WithNewStep("Подготовка: Наполнение корзины", func(t provider.StepCtx) {
			request, err := getAddItemRequest(s.Host, addItemRequest, userID, sku)
			require.NoError(t, err)

			response, err := s.client.Do(request)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, response.StatusCode)

			request2, err := getAddItemRequest(s.Host, addItemRequest2, userID, sku2)
			require.NoError(t, err)

			response2, err := s.client.Do(request2)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, response2.StatusCode)
		})
	})

	t.WithNewStep("Действие: Получение", func(t provider.StepCtx) {

		request, err := getGetCartRequest(s.Host, userID)
		require.NoError(t, err)

		response, err := s.client.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		reportCart, err := decodeResponseBody(response)
		require.NoError(t, err)

		sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })

		require.Equal(t, 2, len(reportCart.Items))
		require.Equal(t, sku, reportCart.Items[0].SKU)
		require.Equal(t, sku2, reportCart.Items[1].SKU)
		require.Equal(t, countUint32, reportCart.Items[0].Count)
		require.Equal(t, count2Uint32, reportCart.Items[1].Count)

	})

	t.WithNewStep("Действие: удаление sku из корзины", func(t provider.StepCtx) {
		request, err := getDeleteItemRequest(s.Host, userID, sku)
		require.NoError(t, err)

		response, err := s.client.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, response.StatusCode)
	})

	t.WithNewStep("Проверка: Проверка корзины после удаления", func(t provider.StepCtx) {
		request, err := getGetCartRequest(s.Host, userID)
		require.NoError(t, err)

		response, err := s.client.Do(request)
		require.NoError(t, err)

		reportCart, err := decodeResponseBody(response)
		require.NoError(t, err)

		require.Equal(t, 1, len(reportCart.Items))
	})

	t.WithNewStep("Добавление несуществующего SKU в корзину", func(sCtx provider.StepCtx) {
		//statusCode := s.cartClient.AddItem(ctx, sCtx, userID, int64(wrongSku), 1)

		data, err := json.Marshal(testAddItemRequest{
			Count: count,
		})
		t.Require().NoError(err, "сериализация тела запроса")
		t.WithNewAttachment("request payload", allure.JSON, data)

		url := fmt.Sprintf("%s/user/%d/cart/%d", s.Host, userID, wrongSku)
		r, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, bytes.NewBuffer(data))
		t.Require().NoError(err, "создание запроса")

		r.Header.Add("Content-Type", "application/json")

		res, err := s.client.Do(r)
		t.Require().NoError(err, "выполнение запроса")
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		t.Require().NoError(err, "не удалось считать ответ")
		t.WithNewAttachment("response body", allure.JSON, body)

		StatusCode(sCtx, http.StatusPreconditionFailed, res.StatusCode)
	})
}

func getDeleteItemRequest(host string, userID int64, sku int64) (*http.Request, error) {
	request, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/user/%s/cart/%s", host, strconv.FormatInt(userID, 10), strconv.FormatInt(sku, 10)),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

func getAddItemRequest(host string, addItemRequest testAddItemRequest, userID int64, sku int64) (*http.Request, error) {
	body, err := json.Marshal(addItemRequest)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(body)
	request, err := http.NewRequest(
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

func getGetCartRequest(host string, userID int64) (*http.Request, error) {
	request, err := http.NewRequest(
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

func getClearCartRequest(host string, userID int64) (*http.Request, error) {
	request, err := http.NewRequest(
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

func StatusCode(sCtx provider.StepCtx, expected, actual int) {
	sCtx.Require().Equal(expected, actual, "Не совпадает статус код")
}
