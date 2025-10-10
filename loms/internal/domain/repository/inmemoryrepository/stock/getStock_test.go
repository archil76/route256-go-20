package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_GetStock_Table(t *testing.T) {
	testData := []testDataStruct{
		{
			name: "Успешное получение стока",
			sku:  tp.sku,
			stock: model.Stock{
				Sku:        tp.sku,
				TotalCount: 10,
				Reserved:   1,
			},
			wantedErr: nil,
		},
		{
			name: "Успешное получение стока 2",
			sku:  tp.sku2,
			stock: model.Stock{
				Sku:        tp.sku2,
				TotalCount: 10,
				Reserved:   1,
			},
			wantedErr: nil,
		},
		{
			name: "Не успешное получение стока с отрицательным ID",
			sku:  0,

			wantedErr: model.ErrSkuIsNotValid,
		},
		{
			name: "Не успешное получение валидно пустого стока",
			sku:  tp.sku3,

			wantedErr: model.ErrStockDoesntExist,
		},
	}

	ctx := context.Background()

	handler := NewStockInMemoryRepository(10)

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			_, err := handler.GetStock(ctx, td.sku)
			require.ErrorIs(t, err, td.wantedErr)

		})
	}
}
