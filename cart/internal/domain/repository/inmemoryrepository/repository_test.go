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
			userID:    tp.wrongUserIDZero,
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
