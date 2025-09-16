package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	desc "route256/loms/internal/api"

	"github.com/ozontech/allure-go/pkg/framework/asserts_wrapper/require"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerE struct {
	suite.Suite

	LomsHostGRPS string
	LomsHostHTTP string
	CartHost     string

	clientHTTP *http.Client
	clientGRPS desc.LomsClient
}

func TestServerE(t *testing.T) {
	t.Parallel()
	suite.RunSuite(t, new(ServerE))
}

func (s *ServerE) BeforeAll(t provider.T) {
	s.CartHost = "http://localhost:8080"
	s.LomsHostHTTP = "http://localhost:8084"
	s.LomsHostGRPS = "http://localhost:8083"

	s.clientHTTP = &http.Client{}

	conn, err := grpc.NewClient(
		":8083",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	s.clientGRPS = desc.NewLomsClient(conn)

	t.Title("loms e2e test")
}

type testItem struct {
	SKU   int64
	Count uint32
}

type testCreateOrderRequest struct {
	userID int64
	items  []testItem
}

func (s *ServerE) TestServerE(t provider.T) {
	t.Parallel()
	t.Helper()

	sku := int64(1076963)
	sku2 := int64(135717466) // должен быть больше sku для проверки сортировки получаемой корзины

	count := uint32(3)
	count2 := uint32(2)
	userID := int64(31337)

	//countUint32 := uint32(count)   //nolint:gosec
	//count2Uint32 := uint32(count2) //nolint:gosec

	createOrderRequest := testCreateOrderRequest{
		userID: userID,
		items: []testItem{
			{
				SKU:   sku,
				Count: count,
			},
			{
				SKU:   sku2,
				Count: count2,
			},
		},
	}

	t.Title("LOMs")

	t.WithNewStep("create normal order", func(t provider.StepCtx) {
		request, err := getCreateOrderRequest(s.LomsHostHTTP, createOrderRequest)
		require.NoError(t, err)

		response, err := s.clientHTTP.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	//t.WithNewStep("Подготовка: Проверка что корзина пуста", func(t provider.StepCtx) {
	//	request, err := getGetCartRequest(s.Host, userID)
	//	require.NoError(t, err)
	//
	//	response, err := s.client.Do(request)
	//	require.NoError(t, err)
	//
	//	require.Equal(t, http.StatusNotFound, response.StatusCode)
	//})
	//
	//t.WithNewStep("Подготовка: Наполнение корзины", func(t provider.StepCtx) {
	//	request, err := getAddItemRequest(s.Host, addItemRequest, userID, sku)
	//	require.NoError(t, err)
	//
	//	response, err := s.client.Do(request)
	//	require.NoError(t, err)
	//	require.Equal(t, http.StatusOK, response.StatusCode)
	//
	//	request2, err := getAddItemRequest(s.Host, addItemRequest2, userID, sku2)
	//	require.NoError(t, err)
	//
	//	response2, err := s.client.Do(request2)
	//	require.NoError(t, err)
	//	require.Equal(t, http.StatusOK, response2.StatusCode)
	//})
	//
	//t.WithNewStep("Действие: Получение", func(t provider.StepCtx) {
	//
	//	request, err := getGetCartRequest(s.Host, userID)
	//	require.NoError(t, err)
	//
	//	response, err := s.client.Do(request)
	//	require.NoError(t, err)
	//	require.Equal(t, http.StatusOK, response.StatusCode)
	//
	//	reportCart, err := decodeResponseBody(response)
	//	require.NoError(t, err)
	//
	//	sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })
	//
	//	require.Equal(t, 2, len(reportCart.Items))
	//	require.Equal(t, sku, reportCart.Items[0].SKU)
	//	require.Equal(t, sku2, reportCart.Items[1].SKU)
	//	require.Equal(t, countUint32, reportCart.Items[0].Count)
	//	require.Equal(t, count2Uint32, reportCart.Items[1].Count)
	//
	//})
	//
	//t.WithNewStep("Действие: удаление sku из корзины", func(t provider.StepCtx) {
	//	request, err := getDeleteItemRequest(s.Host, userID, sku)
	//	require.NoError(t, err)
	//
	//	response, err := s.client.Do(request)
	//	require.NoError(t, err)
	//	require.Equal(t, http.StatusNoContent, response.StatusCode)
	//})
	//
	//t.WithNewStep("Проверка: Проверка корзины после удаления", func(t provider.StepCtx) {
	//	request, err := getGetCartRequest(s.Host, userID)
	//	require.NoError(t, err)
	//
	//	response, err := s.client.Do(request)
	//	require.NoError(t, err)
	//
	//	reportCart, err := decodeResponseBody(response)
	//	require.NoError(t, err)
	//
	//	require.Equal(t, 1, len(reportCart.Items))
	//})
}
func getCreateOrderRequest(host string, createOrderRequest testCreateOrderRequest) (*http.Request, error) {
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

//func getDeleteItemRequest(host string, userID int64, sku int64) (*http.Request, error) {
//	request, err := http.NewRequest(
//		http.MethodDelete,
//		fmt.Sprintf("%s/user/%s/cart/%s", host, strconv.FormatInt(userID, 10), strconv.FormatInt(sku, 10)),
//		http.NoBody,
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	request.Header.Add("Content-Type", "application/json")
//
//	return request, nil
//}
//
//
//
//func getGetCartRequest(host string, userID int64) (*http.Request, error) {
//	request, err := http.NewRequest(
//		http.MethodGet,
//		fmt.Sprintf("%s/user/%s/cart", host, strconv.FormatInt(userID, 10)),
//		http.NoBody,
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	request.Header.Add("Content-Type", "application/json")
//
//	return request, nil
//}
//
//func getClearCartRequest(host string, userID int64) (*http.Request, error) {
//	request, err := http.NewRequest(
//		http.MethodDelete,
//		fmt.Sprintf("%s/user/%s/cart", host, strconv.FormatInt(userID, 10)),
//		http.NoBody,
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	request.Header.Add("Content-Type", "application/json")
//
//	return request, nil
//}
//
//func decodeResponseBody(response *http.Response) (testReportCart, error) {
//	defer response.Body.Close()
//
//	decoder := json.NewDecoder(response.Body)
//
//	reportCart := testReportCart{}
//
//	err := decoder.Decode(&reportCart)
//
//	return reportCart, err
//}
