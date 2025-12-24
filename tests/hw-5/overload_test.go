package hw5

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/clients/cart"
)

const (
	maxWorkerPoolJobs = 11
	totalRequests     = 70
	totalBurstsMin    = 2
	totalBurstsMax    = 3
	reqsPerBurstMin   = 30
	reqsPerBurstMax   = 50
)

func worker(
	ctx context.Context,
	sCtx provider.StepCtx,
	jobs <-chan struct{},
	statusCh chan<- int,
	wg *sync.WaitGroup,
	client *cart.Client,
	userID int64,
) {
	for range jobs {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, code := client.GetCart(ctx, sCtx, userID)
		statusCh <- code
		wg.Done()
	}
}

func (s *Suite) TestWithSoftOverload(t provider.T) {
	err := s.setToxicEnabled(s.toxicName, false)
	t.Require().NoError(err, "latency disabled")
	t.Title("Нагрузка с равномерным трафиком")

	userID := s.randomize.Int64()
	t.WithNewParameters("userID", userID)

	t.WithNewStep("Подготовка корзины", func(sCtx provider.StepCtx) {
		sCtx.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
			code := s.cartClient.DeleteCart(s.ctx, sCtx, userID)
			sCtx.Assert().Equal(http.StatusNoContent, code)
		})
		sCtx.WithNewStep("Ожидается пустая корзина", func(sCtx provider.StepCtx) {
			getCart, code := s.cartClient.GetCart(s.ctx, sCtx, userID)
			sCtx.Assert().Equal(http.StatusNotFound, code)
			sCtx.Assert().Empty(getCart.Items)
		})
		sCtx.WithNewStep("Наполнение корзины", func(sCtx provider.StepCtx) {
			for _, sku := range s.totalSkus {
				code := s.cartClient.AddItem(s.ctx, sCtx, userID, sku, 10)
				sCtx.Assert().Equal(http.StatusOK, code)
			}
		})
	})

	t.WithNewStep("Нагрузка > 10 RPS", func(sCtx provider.StepCtx) {
		ctx, cancel := context.WithCancel(s.ctx)
		defer cancel()

		var (
			wg       sync.WaitGroup
			jobs     = make(chan struct{}, maxWorkerPoolJobs)
			statusCh = make(chan int, totalRequests)
			counts   = make(map[int]int)
		)

		wg.Add(totalRequests)
		for i := 0; i < maxWorkerPoolJobs; i++ {
			go worker(ctx, sCtx, jobs, statusCh, &wg, s.cartClient, userID)
		}

		ticker := time.NewTicker(1 * time.Millisecond)
		defer ticker.Stop()

		for i := 0; i < totalRequests; i++ {
			<-ticker.C
			jobs <- struct{}{}
		}
		close(jobs)

		wg.Wait()
		close(statusCh)

		for code := range statusCh {
			counts[code]++
		}

		for code, cnt := range counts {
			sCtx.WithNewParameters(fmt.Sprintf("Status %d", code), cnt)
		}

		successful := counts[http.StatusOK]
		tooMany := counts[http.StatusTooManyRequests]
		serverErrors := 0
		for code, cnt := range counts {
			if code != http.StatusOK && code != http.StatusTooManyRequests && code != http.StatusNoContent {
				serverErrors += cnt
			}
		}

		sCtx.Require().Equal(0, tooMany, "Сервис не должен отвечать 429")
		sCtx.Require().Equal(0, serverErrors, "Сервис не должен отвечать 50x")
		sCtx.Require().Equal(totalRequests, successful+tooMany+serverErrors, "Общее количество запросов должно совпадать")

		sCtx.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
			code := s.cartClient.DeleteCart(s.ctx, sCtx, userID)
			sCtx.Assert().Equal(http.StatusNoContent, code)
		})
		sCtx.WithNewStep("Ожидается пустая корзина", func(sCtx provider.StepCtx) {
			getCart, code := s.cartClient.GetCart(s.ctx, sCtx, userID)
			sCtx.Assert().Equal(http.StatusNotFound, code)
			sCtx.Assert().Empty(getCart.Items)
		})
	})
}

func (s *Suite) TestRateLimitWithBurstLoad(t provider.T) {
	err := s.setToxicEnabled(s.toxicName, false)
	t.Require().NoError(err, "latency disabled")

	t.Title("Нагрузка с всплесками трафика")

	userID := s.randomize.Int64()
	t.WithNewParameters("userID", userID)

	t.WithNewStep("Подготовка корзины", func(sCtx provider.StepCtx) {
		sCtx.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
			code := s.cartClient.DeleteCart(s.ctx, sCtx, userID)
			sCtx.Assert().Equal(http.StatusNoContent, code)
		})
		sCtx.WithNewStep("Ожидается пустая корзина", func(sCtx provider.StepCtx) {
			cart, code := s.cartClient.GetCart(s.ctx, sCtx, userID)
			sCtx.Assert().Equal(http.StatusNotFound, code)
			sCtx.Assert().Empty(cart.Items)
		})
		sCtx.WithNewStep("Наполнение корзины", func(sCtx provider.StepCtx) {
			for _, sk := range s.totalSkus {
				code := s.cartClient.AddItem(s.ctx, sCtx, userID, sk, 10)
				sCtx.Assert().Equal(http.StatusOK, code)
			}
		})
	})

	burstsCount := s.randomize.IntN(totalBurstsMax-totalBurstsMin+1) + totalBurstsMin
	requestsPerBurst := s.randomize.IntN(reqsPerBurstMax-reqsPerBurstMin+1) + reqsPerBurstMin
	pauseBetweenBursts := time.Duration(s.randomize.IntN(200)+100) * time.Millisecond

	t.WithNewStep("Отправляем трафик всплесками", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			"Количество всплесков", burstsCount,
			"RPS на всплеск", requestsPerBurst,
			"Пауза между всплесками (ms)", pauseBetweenBursts.Milliseconds(),
		)

		ctx, cancel := context.WithCancel(s.ctx)
		defer cancel()

		totalJobs := burstsCount * requestsPerBurst
		var (
			wg       sync.WaitGroup
			jobs     = make(chan struct{}, totalJobs)
			statusCh = make(chan int, totalJobs)
			counts   = make(map[int]int)
		)

		for i := 0; i < requestsPerBurst; i++ {
			go worker(ctx, sCtx, jobs, statusCh, &wg, s.cartClient, userID)
		}

		wg.Add(totalJobs)
		for burst := 1; burst <= burstsCount; burst++ {
			sCtx.Logf("Всплеск %d: %d запросов", burst, requestsPerBurst)
			for i := 0; i < requestsPerBurst; i++ {
				jobs <- struct{}{}
			}
			if burst < burstsCount {
				time.Sleep(pauseBetweenBursts)
			}
		}
		close(jobs)

		wg.Wait()
		close(statusCh)

		for code := range statusCh {
			counts[code]++
		}

		for code, cnt := range counts {
			sCtx.WithNewParameters(fmt.Sprintf("Status %d", code), cnt)
		}

		successful := counts[http.StatusOK]
		rateLimited := counts[http.StatusTooManyRequests]
		serverErrors := 0
		for code, cnt := range counts {
			if code != http.StatusOK && code != http.StatusTooManyRequests {
				serverErrors += cnt
			}
		}

		sCtx.Require().Equal(totalJobs, successful+rateLimited+serverErrors, "Общее количество запросов")
		sCtx.Require().Equal(0, rateLimited, "Сервис не должен отвечать 429")
		sCtx.Require().Equal(0, serverErrors, "Сервис не должен отвечать 50x")

		sCtx.WithNewStep("Очистка корзины", func(sCtx provider.StepCtx) {
			code := s.cartClient.DeleteCart(s.ctx, sCtx, userID)
			sCtx.Assert().Equal(http.StatusNoContent, code)
		})
		sCtx.WithNewStep("Проверка пустой корзины", func(sCtx provider.StepCtx) {
			userCart, code := s.cartClient.GetCart(s.ctx, sCtx, userID)

			sCtx.Assert().Equal(http.StatusNotFound, code)
			sCtx.Assert().Empty(userCart.Items)
		})
	})
}
