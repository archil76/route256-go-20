package inmemoryrepository

import (
	"context"
	"route256/cart/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

type testDataStruct struct {
	name      string
	userID    int64
	item      model.Item
	wantedErr error
}

var (
	tp = struct {
		sku             model.Sku
		sku2            model.Sku
		userID          int64
		userID2         int64
		userID3         int64
		wrongSkuZero    model.Sku
		wrongUserIDZero int64
		wrongUserIDNeg  int64
		count           uint32
		count2          uint32
	}{
		sku:             model.Sku(100),
		sku2:            model.Sku(101),
		userID:          123,
		userID2:         2525455,
		userID3:         1001,
		wrongSkuZero:    model.Sku(0),
		wrongUserIDZero: 0,
		wrongUserIDNeg:  -125555,
		count:           2,
		count2:          3,
	}
)

func TestHandler_All(t *testing.T) {
	t.Run("TestHandler_AddItem_Table", TestHandler_AddItem_Table)
	t.Run("TestHandler_CreateCart_Table", TestHandler_CreateCart_Table)
	t.Run("TestHandler_GetCart_Table", TestHandler_GetCart_Table)
	t.Run("TestHandler_DeleteItem_Table", TestHandler_DeleteItem_Table)
	t.Run("TestHandler_DeleteItems_Table", TestHandler_DeleteItems_Table)
}

// наполняет тестовыми данными корзину перед тестированием методов получения и удаления
func fillCart(ctx context.Context, t *testing.T, addHandler *Repository) map[model.UserID]map[model.Sku]uint32 {
	testAddData := map[model.UserID]map[model.Sku]uint32{
		tp.userID: {
			tp.sku:  tp.count,
			tp.sku2: tp.count2,
		},
		tp.userID2: {
			tp.sku: tp.count,
		},
		tp.userID3: {},
	}

	for key, value := range testAddData {
		for sku, count := range value {
			_, err := addHandler.AddItem(ctx, key, model.Item{Sku: sku, Count: count})
			require.ErrorIs(t, err, nil)
		}
	}

	return testAddData
}
