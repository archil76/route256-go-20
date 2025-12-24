package hw5

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (s *Suite) TestDataRace(t provider.T) {
	err := s.setToxicEnabled(s.toxicName, false)
	t.Require().NoError(err, "latency disabled")

	t.Title("Параллельное создание заказов")

	userID := s.randomize.Int64()

	t.WithNewParameters("userID", userID)
	t.WithNewStep("Подготовка корзины", func(sCtx provider.StepCtx) {
		sCtx.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
			code := s.cartClient.DeleteCart(s.ctx, sCtx, userID)
			sCtx.Assert().Equal(http.StatusNoContent, code)
		})
		sCtx.WithNewStep("Проверка пустой корзины", func(sCtx provider.StepCtx) {
			cart, code := s.cartClient.GetCart(s.ctx, sCtx, userID)
			sCtx.Assert().Equal(http.StatusNotFound, code)
			sCtx.Assert().Empty(cart.Items)
		})
	})

	t.WithNewStep("Создание заказов параллельно", func(sCtx provider.StepCtx) {
		ctx, cancel := context.WithCancel(s.ctx)
		defer cancel()

		totalJobs := len(s.totalSkus)

		var wg sync.WaitGroup
		results := make(chan struct{ code, sku int }, totalJobs)

		wg.Add(totalJobs)
		for _, sku := range s.totalSkus {
			go func(sku int64) {
				defer wg.Done()
				code := s.cartClient.AddItem(ctx, sCtx, userID, sku, 10)
				results <- struct{ code, sku int }{code, int(sku)}
			}(sku)
		}

		wg.Wait()
		close(results)

		counts := make(map[int]int)
		// Логируем коды и SKU после всех горутин
		for r := range results {
			sCtx.Logf("Код %d, sku %d", r.code, r.sku)
			counts[r.code]++
		}

		for code, cnt := range counts {
			sCtx.WithNewParameters(fmt.Sprintf("Status %d", code), cnt)
		}

		sCtx.Require().Equal(totalJobs, counts[http.StatusOK], "Все заказы должны успешно создаваться")
	})

	t.WithNewStep("Проверка состава корзины", func(sCtx provider.StepCtx) {
		cart, code := s.cartClient.GetCart(s.ctx, sCtx, userID)

		sCtx.Assert().Equal(http.StatusOK, code)
		sCtx.Require().NotEmpty(cart.Items, "Корзина не должна быть пуста")
		sCtx.Require().True(checkSkusMatch(s.totalSkus, cart.Items), "Количество товаров в корзине должно быть консистентным")
	})

	t.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
		statusCode := s.cartClient.DeleteCart(s.ctx, sCtx, userID)

		sCtx.Assert().Equal(http.StatusNoContent, statusCode)
	})

	t.WithNewStep("Проверка пустой корзины", func(sCtx provider.StepCtx) {
		cart, statusCode := s.cartClient.GetCart(s.ctx, sCtx, userID)

		sCtx.Assert().Equal(http.StatusNotFound, statusCode)
		sCtx.Assert().Empty(cart.Items)
	})
}
