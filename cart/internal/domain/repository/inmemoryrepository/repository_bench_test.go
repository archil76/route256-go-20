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
		go func(i int) {
			wg.Add(1)
			userID := int64(rand.Int63() + 1)
			sku := model.Sku(rand.Int63() + 1)
			count := uint32(i + 100/10)
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
