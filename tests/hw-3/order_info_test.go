package hw3

import (
	"context"
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
)

func (s *Suite) TestOrderInfo_InvalidRequest(t provider.T) {
	t.Title("Неуспешное получение заказа из-за невалидного запроса")

	var (
		ctx = context.Background()
	)

	t.Run("Нулевой идентификатор заказа", func(t provider.T) {
		t.WithNewStep("Получение заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, 0)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})

	t.Run("Отрицательный идентификатор заказа", func(t provider.T) {
		t.WithNewStep("Получение заказа", func(sCtx provider.StepCtx) {
			_, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, -1)
			assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		})
	})
}

func (s *Suite) TestOrderInfo_NotExistingOrder(t provider.T) {
	t.Title("Получение несуществующего заказа")

	var (
		ctx = context.Background()
	)

	t.WithNewStep("Получение заказа", func(sCtx provider.StepCtx) {
		_, statusCode := s.lomsClient.OrderInfo(ctx, sCtx, 9223372036854775001)
		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
	})
}
