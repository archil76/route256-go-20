package inmemoryrepository

import (
	"context"
	"math/rand"
	"route256/cart/internal/domain/model"
	"testing"
)

type TestCase struct {
	userID int64
	item   *model.Item
}

func BenchmarkRepository_AddItem(b *testing.B) {
	ctx := context.Background()

	handler := NewCartInMemoryRepository(100)

	testData := make([]TestCase, b.N)
	for n := 0; n < b.N; n++ {
		userID := rand.Int63() + 1    //nolint:gosec
		sku := rand.Int63() + 1       //nolint:gosec
		count := uint32(rand.Int31()) //nolint:gosec
		item := model.Item{Sku: sku, Count: count}

		testData[n] = TestCase{
			userID: userID,
			item:   &item,
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {

		var i = -1
		for pb.Next() {
			i++
			testCase := testData[i]
			_, err := handler.AddItem(ctx, testCase.userID, *testCase.item)
			if err != nil {

			}
		}
	})
}

func BenchmarkRepository_DeleteItem(b *testing.B) {
	ctx := context.Background()

	handler := NewCartInMemoryRepository(100)

	testData := make([]TestCase, b.N)
	for n := 0; n < b.N; n++ {
		userID := rand.Int63() + 1    //nolint:gosec
		sku := rand.Int63() + 1       //nolint:gosec
		count := uint32(rand.Int31()) //nolint:gosec
		item := model.Item{Sku: sku, Count: count}

		testData[n] = TestCase{
			userID: userID,
			item:   &item,
		}

		_, err := handler.AddItem(ctx, userID, item)
		if err != nil {
			b.Fatalf("Плохой пример с параметрами userID: %d; SKU: %d; Count: %d\n", userID, item.Sku, item.Count)
			return
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i = -1
		for pb.Next() {
			i++
			testCase := testData[i]
			_, err := handler.DeleteItem(ctx, testCase.userID, *testCase.item)
			if err != nil {

			}
		}
	})
}
