package inmemoryrepository

import (
	"context"
	"fmt"
	"math/rand"
	"route256/cart/internal/domain/model"
	"sync"
	"testing"
)

func BenchmarkRepository_AddItem(b *testing.B) {
	wg := sync.WaitGroup{}
	ctx := context.Background()

	handler := NewCartInMemoryRepository(100)

	for i := 0; i < b.N; i++ {
		wg.Add(1)

		go func(i int) {
			userID := rand.Int63() + 1  //nolint:gosec
			sku := rand.Int63() + 1     //nolint:gosec
			count := uint32(i + 100/10) //nolint:gosec
			item := model.Item{Sku: sku, Count: count}

			_, err := handler.AddItem(ctx, userID, item)

			if err != nil {
				fmt.Printf("err %v\n", err)
			}

			_, err = handler.DeleteItem(ctx, userID, item)
			if err != nil {
				fmt.Printf("err %v\n", err)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
