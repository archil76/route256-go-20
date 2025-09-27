package inmemoryrepository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_GetCart_Table(t *testing.T) {
	testData := []testDataStruct{
		{
			name:      "Успешное получение корзины",
			userID:    tp.userID,
			wantedErr: nil,
		},
		{
			name:      "Успешное получение корзины 2",
			userID:    tp.userID2,
			wantedErr: nil,
		},
		{
			name:      "Не успешное получение корзины пользователя с отрицательным ID",
			userID:    tp.wrongUserIDNeg,
			wantedErr: ErrUserIDIsNotValid,
		},
		{
			name:      "Не успешное получение Пустой корзины валидного пользователя",
			userID:    tp.userID3,
			wantedErr: ErrCartDoesntExist,
		},
	}

	ctx := context.Background()

	addHandler := NewCartInMemoryRepository(5)

	testAddData := fillCart(ctx, t, addHandler)

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {

			cart, err := addHandler.GetCart(ctx, td.userID)
			require.ErrorIs(t, err, td.wantedErr)
			if err == nil {
				require.Equal(t, td.userID, cart.UserID)
				require.Equal(t, len(testAddData[cart.UserID]), len(cart.Items))

				gotCart := testAddData[cart.UserID]
				for sku, count := range cart.Items {
					require.Equal(t, count, gotCart[sku])
				}
			}
		})
	}
}
