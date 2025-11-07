//go:build e2e_test

package e2e

import (
	"context"
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/asserts_wrapper/require"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *ServerE) TestServerParallelWrongSku(t provider.T) {

	t.Helper()

	userID := int64(1022222)
	skus := []int64{
		32638658,
		32605854,
		32205848,
		32205849,
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
