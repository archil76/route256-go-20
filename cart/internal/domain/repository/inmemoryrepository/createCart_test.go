package inmemoryrepository

import (
	"context"
	"route256/cart/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_CreateCart_Table(t *testing.T) {
	testData := []testDataStruct{
		{
			name:      "Успешное создание корзины",
			userID:    tp.userID,
			wantedErr: nil,
		},
		{
			name:      "Успешное создание корзины 2",
			userID:    tp.userID2,
			wantedErr: nil,
		},
		{
			name:      "Не падает в ошибку при создании корзины если она уже есть.",
			userID:    tp.userID,
			wantedErr: nil,
		},
		{
			name:      "Не успешное создании корзины пользователю с нулевым ID",
			userID:    tp.wrongUserIDZero,
			wantedErr: ErrUserIDIsNotValid,
		},
		{
			name:      "Не успешное создании корзины пользователю с отрицательным ID",
			userID:    tp.wrongUserIDNeg,
			wantedErr: ErrUserIDIsNotValid,
		},
	}

	ctx := context.Background()

	addHandler := NewCartInMemoryRepository(100)

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			cart := model.Cart{UserID: td.userID, Items: map[model.Sku]uint32{}}

			newCart, err := addHandler.createCart(ctx, cart)

			require.ErrorIs(t, err, td.wantedErr)

			if td.wantedErr == nil {
				require.Equal(t, newCart.UserID, cart.UserID)
			}
		})
	}
}
