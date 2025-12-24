package hw5

import (
	"context"
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestGetListSuccess(t provider.T) {
	err := s.setToxicEnabled(s.toxicName, true)
	t.Require().NoError(err, "latency enabled")

	t.Title("Получение корзины")

	userID := s.randomize.Int64()
	ctx, cancel := context.WithTimeout(s.ctx, 11*time.Second)
	defer cancel()

	t.WithNewParameters("userID", userID)

	t.WithNewStep("Подготовка корзины", func(sCtx provider.StepCtx) {
		sCtx.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
			statusCode := s.cartClient.DeleteCart(ctx, sCtx, userID)

			sCtx.Assert().Equal(http.StatusNoContent, statusCode)
		})

		sCtx.WithNewStep("Ожидается пустая корзина", func(sCtx provider.StepCtx) {
			cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

			sCtx.Assert().Equal(http.StatusNotFound, statusCode)
			sCtx.Assert().Empty(cart.Items)
		})

		sCtx.WithNewStep("Наполнение корзины", func(sCtx provider.StepCtx) {
			for _, sku := range s.totalSkus {
				statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku, 10)
				sCtx.Assert().Equal(http.StatusOK, statusCode)
			}
		})
	})

	t.WithNewStep("Получение списка заказов", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

		sCtx.Assert().Equal(http.StatusOK, statusCode)
		sCtx.Assert().NotEmpty(cart.Items)
		sCtx.Require().True(checkSkusMatch(s.totalSkus, cart.Items), "Количество товаров в корзине должно быть консистентным")

		t.LogStep("cart", cart)
	})

	t.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.DeleteCart(ctx, sCtx, userID)

		sCtx.Assert().Equal(http.StatusNoContent, statusCode)
	})

	t.WithNewStep("Проверка пустой корзины", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

		sCtx.Assert().Equal(http.StatusNotFound, statusCode)
		sCtx.Assert().Empty(cart.Items)
	})
}
