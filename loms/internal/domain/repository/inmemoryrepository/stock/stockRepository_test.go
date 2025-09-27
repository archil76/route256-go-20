package inmemoryrepository

import (
	"route256/loms/internal/domain/model"
	"testing"
)

type testDataStruct struct {
	name      string
	sku       int64
	stock     model.Stock
	wantedErr error
}

var (
	tp = struct {
		sku  int64
		sku2 int64
		sku3 int64
	}{
		sku:  139275865,
		sku2: 2956315,
		sku3: 1001,
	}
)

func TestHandler_All(t *testing.T) {
	t.Run("TestHandler_GetStock_Table", TestHandler_GetStock_Table)
	//t.Run("TestHandler_CreateCart_Table", TestHandler_CreateCart_Table)
	//t.Run("TestHandler_GetCart_Table", TestHandler_GetCart_Table)
	//t.Run("TestHandler_DeleteItem_Table", TestHandler_DeleteItem_Table)
	//t.Run("TestHandler_DeleteItems_Table", TestHandler_DeleteItems_Table)
}
