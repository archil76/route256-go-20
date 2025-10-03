package service

import (
	orderrepo "route256/loms/internal/domain/repository/inmemoryrepository/order"
	stockrepo "route256/loms/internal/domain/repository/inmemoryrepository/stock"
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

func NewLomsServiceWithInMemoryRepository() *LomsService {
	orderRepository := orderrepo.NewOrderInMemoryRepository(10, &counter)

	stockRepository := stockrepo.NewStockInMemoryRepository(10)

	return NewLomsService(orderRepository, stockRepository, nil)
}

func TestHandler_All(t *testing.T) {
	t.Run("Test_OrderCancel", Test_OrderCancel)
	t.Run("Test_OrderCreate", Test_OrderCreate)
	t.Run("Test_OrderInfo", Test_OrderInfo)
	t.Run("Test_OrderPay", Test_OrderPay)
	t.Run("Test_StockInfo", Test_StockInfo)
}
