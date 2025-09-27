package inmemoryrepository

import (
	"context"
	"route256/cart/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_AddItem_Table(t *testing.T) {
	testData := []testDataStruct{
		{
			name:   "Успешное добавление товара",
			userID: tp.userID,
			item: model.Item{
				Sku:   tp.sku,
				Count: tp.count,
			},
			wantedErr: nil,
		},
		{
			name:   "Успешное добавление товара",
			userID: tp.userID,
			item: model.Item{
				Sku:   tp.sku,
				Count: tp.count2,
			},
			wantedErr: nil,
		},
		{
			name:   "Успешное добавление товара 2",
			userID: tp.userID,
			item: model.Item{
				Sku:   tp.sku2,
				Count: tp.count2,
			},
			wantedErr: nil,
		},
		{
			name:   "Не успешное добавление товара пользователю с нулевым ID",
			userID: tp.wrongUserIDZero,
			item: model.Item{
				Sku:   model.Sku(tp.sku2),
				Count: tp.count2,
			},
			wantedErr: ErrUserIDIsNotValid,
		},
	}

	ctx := context.Background()

	addHandler := NewCartInMemoryRepository(100)

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			item, err := addHandler.AddItem(ctx, td.userID, td.item)

			require.ErrorIs(t, err, td.wantedErr)

			if td.wantedErr == nil {
				require.Equal(t, td.item.Sku, item.Sku)
				require.Equal(t, td.item.Count, item.Count)
			}
		})
	}
}
