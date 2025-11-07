//go:build e2e_test

package e2e

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/asserts_wrapper/require"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

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

	ctx, cancel := context.WithTimeout(s.ctx, 11*time.Second)
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
