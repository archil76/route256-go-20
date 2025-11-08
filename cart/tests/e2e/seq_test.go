//go:build e2e_test

package e2e

import (
	"net/http"
	"sort"

	"github.com/ozontech/allure-go/pkg/framework/asserts_wrapper/require"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestServer_Seq(t provider.T) {
	t.Helper()

	sku := int64(1076963)
	sku2 := int64(1148162) // должен быть больше sku для проверки сортировки получаемой корзины
	//wrongSku := uint32(1076963000)

	count := int32(2)
	count2 := int32(3)
	userID := int64(1022202)

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
}
