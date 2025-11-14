package hw1

import (
	"context"
	"math"
	"math/rand"
	"net/http"

	"route256/tests/app/assert"

	"route256/tests/app/domain"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestAddSkuSuccess(t provider.T) {
	t.Parallel()

	t.Title("Успешное добавление SKU в корзину")

	var (
		sku      = int64(1076963)
		sku2     = int64(1148162)
		wrongSku = int64(1076963000)
		userID   = rand.Int63()
		ctx      = context.Background()
		maxCount = int64(math.MaxInt64)
	)

	t.WithNewParameters(
		"userID", userID,
		"Основной SKU", sku,
		"Дополнительный SKU", sku2,
		"Несуществующий SKU", wrongSku,
		"Предельное значение Count", maxCount,
	)

	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
			statusCode := s.cartClient.DeleteCart(ctx, sCtx, userID)

			assert.StatusCode(sCtx, http.StatusNoContent, statusCode)
		})

		t.WithNewStep("Проверка что корзина пуста", func(sCtx provider.StepCtx) {
			cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

			assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
			assert.EmptyCart(sCtx, cart)
		})
	})

	t.WithNewStep("Добавление 1 SKU в корзину", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku, 1)

		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	var (
		itemName  string
		itemPrice uint32
	)

	t.WithNewStep("Проверка корзины", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.NotEmptyCart(sCtx, cart)
		sCtx.Assert().Len(cart.Items, 1, "Ожидается один SKU в корзине")
		itemName = cart.Items[0].Name
		itemPrice = cart.Items[0].Price
	})

	t.WithNewStep("Добавление 5 аналогичных SKU в корзину", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku, 5)

		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	var lastCart *domain.Cart

	t.WithNewStep("Проверка корзины после добавления SKU", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		expectedCart := &domain.Cart{
			Items: []domain.CartItem{
				{
					SKU:   uint32(sku),
					Count: 6,
					Name:  itemName,
					Price: itemPrice,
				},
			},
			TotalPrice: 6 * itemPrice,
		}
		assert.Cart(sCtx, expectedCart, cart)

		lastCart = cart
	})

	t.WithNewStep("Попытка добавления несуществующего SKU", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, wrongSku, 1)
		assert.StatusCode(sCtx, http.StatusPreconditionFailed, statusCode)
	})

	t.WithNewStep("Проверка корзины после добавления несуществующего SKU", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.Cart(sCtx, lastCart, cart)
	})

	t.WithNewStep("Добавление второго SKU в корзину", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.AddItem(ctx, sCtx, userID, sku2, 1)

		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Проверка корзины после добавления второго SKU", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(ctx, sCtx, userID)

		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		assert.NotEmptyCart(sCtx, cart)
		sCtx.Assert().Len(cart.Items, 2, "Ожидается два SKU в корзине")
	})
}
