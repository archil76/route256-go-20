//go:build e2e_test

package e2e

import (
	"net/http"
	"sort"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/asserts_wrapper/require"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *ServerE) TestServerParallel(t provider.T) {
	t.Parallel()
	t.Helper()

	userID := int64(1022223)

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

			for _, sku := range []int64{2956315, 135717466, 135937324, 1148162} {
				request, err := getAddItemRequest(s.Host, testAddItemRequest{
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
		time.Sleep(2 * time.Second)
		request, err := getGetCartRequest(s.Host, userID)
		require.NoError(t, err)

		response, err := s.client.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		reportCart, err := decodeResponseBody(response)
		require.NoError(t, err)

		sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })

		require.Equal(t, 5, len(reportCart.Items))
	})

}

func (s *ServerE) TestServerParallel2(t provider.T) {
	t.Parallel()
	t.Helper()

	sku := int64(1076963)
	sku2 := int64(1148162) // должен быть больше sku для проверки сортировки получаемой корзины
	//wrongSku := uint32(1076963000)

	count := int32(2)
	count2 := int32(3)
	userID := int64(1022223)

	//countUint32 := uint32(count) //nolint:gosec
	//count2Uint32 := uint32(count2) //nolint:gosec

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
		time.Sleep(2 * time.Second)

		request, err := getGetCartRequest(s.Host, userID)
		require.NoError(t, err)

		response, err := s.client.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		reportCart, err := decodeResponseBody(response)
		require.NoError(t, err)

		sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })

		require.Equal(t, 5, len(reportCart.Items))
	})

	t.WithNewStep("Действие: удаление sku из корзины", func(t provider.StepCtx) {
		time.Sleep(2 * time.Second)

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

		require.Equal(t, 4, len(reportCart.Items))
		//require.NotContains(t, model.Item{
		//	Sku:   sku,
		//	Count: count,
		//}, reportCart.Items)
	})
}
