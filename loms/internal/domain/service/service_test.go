package service

import (
	"sync/atomic"
	"testing"
)

var counter atomic.Int64
var (
	tp = struct {
		sku    int64
		sku2   int64
		sku3   int64
		count  uint32
		count2 uint32
		count3 uint32
	}{
		sku:    139275865,
		sku2:   2956315,
		sku3:   1001,
		count:  2,
		count2: 3,
		count3: 15,
	}
)

func TestHandler_All(t *testing.T) {
	t.Run("Test_OrderCreate", Test_OrderCreate)
	t.Run("Test_StockInfo", Test_StockInfo)
	//t.Run("TestHandler_GetCart_Table", TestHandler_GetCart_Table)
	//t.Run("TestHandler_DeleteItem_Table", TestHandler_DeleteItem_Table)
	//t.Run("TestHandler_DeleteItems_Table", TestHandler_DeleteItems_Table)
}
