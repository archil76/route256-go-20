package inmemoryrepository

import (
	"context"
	"route256/cart/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_DeleteItem_Table(t *testing.T) {
	testData := []testDataStruct{
		{
			name:   "Успешное удаление из корзины",
			userID: tp.userID,
			item: model.Item{
				Sku: tp.sku,
			},
			wantedErr: nil,
		},
		{
			name:   "Успешное удаление из корзины",
			userID: tp.userID,
			item: model.Item{
				Sku: tp.sku2,
			},
			wantedErr: nil,
		},
		{
			name:      "Неуспешное удаление из корзины пользователя с нулевым ID",
			userID:    tp.wrongUserIDZero,
			wantedErr: ErrUserIDIsNotValid,
		},
		{
			name:      "Не успешное удаление из несуществующей корзины валидного пользователя",
			userID:    1002,
			wantedErr: ErrCartDoesntExist,
		},
	}

	ctx := context.Background()

	addHandler := NewCartInMemoryRepository(5)

	_ = fillCart(ctx, t, addHandler)

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {

			_, err := addHandler.DeleteItem(ctx, td.userID, td.item)
			require.ErrorIs(t, err, td.wantedErr)
			if err == nil {
				cart, err := addHandler.GetCart(ctx, td.userID)
				require.ErrorIs(t, err, nil)
				require.NotContains(t, cart.Items, td.item)

			}
		})
	}
}
