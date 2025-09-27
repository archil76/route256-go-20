package inmemoryrepository

import (
	"context"
	"route256/cart/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_DeleteItems_Table(t *testing.T) {
	testData := []testDataStruct{
		{
			name:   "Успешная очистка корзины",
			userID: tp.userID,
			item: model.Item{
				Sku: tp.sku,
			},
			wantedErr: nil,
		},
		{
			name:   "Успешная очистка корзины 2",
			userID: tp.userID,
			item: model.Item{
				Sku: tp.sku2,
			},
			wantedErr: nil,
		},
		{
			name:      "Неуспешная очистка корзины пользователя с нулевым ID",
			userID:    tp.wrongUserIDZero,
			wantedErr: ErrUserIDIsNotValid,
		},
		{
			name:      "Неуспешная очистка несуществующей корзины валидного пользователя",
			userID:    1002,
			wantedErr: ErrCartDoesntExist,
		},
	}

	ctx := context.Background()

	addHandler := NewCartInMemoryRepository(5)

	_ = fillCart(ctx, t, addHandler)

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			_, err := addHandler.DeleteItems(ctx, td.userID)

			require.ErrorIs(t, err, nil)
		})
	}

	for _, td := range testData {
		cart, err := addHandler.GetCart(ctx, td.userID)
		require.ErrorIs(t, err, td.wantedErr)
		if err == nil {
			require.Empty(t, cart.Items)
		}
	}
}
