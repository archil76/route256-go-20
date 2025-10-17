//go:build e2e_test

package e2e

import (
	"net/http"
	"sort"

	"github.com/ozontech/allure-go/pkg/framework/asserts_wrapper/require"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *ServerE) TestServerParallel(t provider.T) {
	t.Parallel()
	t.Helper()

	userID := int64(1022222)

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

		request, err := getGetCartRequest(s.Host, userID)
		require.NoError(t, err)

		response, err := s.client.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		reportCart, err := decodeResponseBody(response)
		require.NoError(t, err)

		sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })

		require.Equal(t, 4, len(reportCart.Items))
	})

}
