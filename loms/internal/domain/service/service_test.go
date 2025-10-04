package service

import (
	"route256/loms/internal/domain/model"
	orderrepo "route256/loms/internal/domain/repository/inmemoryrepository/order"
	stockrepo "route256/loms/internal/domain/repository/inmemoryrepository/stock"
	"route256/loms/internal/domain/service/mock"
	"sync/atomic"
	"testing"

	"github.com/gojuno/minimock/v3"
)

var counter atomic.Int64
var (
	tp = struct {
		orderID, orderID2, userID, userID2, sku, sku2, sku3 int64
		count, count2, count3                               uint32
	}{
		orderID:  1,
		orderID2: 2,
		userID:   12546,
		userID2:  8888,
		sku:      139275865,
		sku2:     2956315,
		sku3:     1001,
		count:    2,
		count2:   3,
		count3:   15,
	}
)

var (
	items = []model.Item{
		{
			Sku:   tp.sku,
			Count: tp.count,
		},
	}

	newOrder = model.Order{
		OrderID: 0,
		UserID:  tp.userID,
		Status:  model.NEWSTATUS,
		Items:   items,
	}
	newOrder2 = model.Order{
		OrderID: 0,
		UserID:  tp.userID2,
		Status:  model.NEWSTATUS,
		Items:   items,
	}
	order = model.Order{
		OrderID: tp.orderID,
		UserID:  tp.userID,
		Status:  model.NEWSTATUS,
		Items:   items,
	}

	order2 = model.Order{
		OrderID: tp.orderID2,
		UserID:  tp.userID2,
		Status:  model.NEWSTATUS,
		Items:   items,
	}

	canceledOrder = model.Order{
		OrderID: tp.orderID,
		UserID:  tp.userID,
		Status:  model.CANCELED,
		Items:   items,
	}
	awaitingPaymentOrder = model.Order{
		OrderID: tp.orderID,
		UserID:  tp.userID,
		Status:  model.AWAITINGPAYMENT,
		Items:   items,
	}
	stocks = []model.Stock{
		{
			Sku:        tp.sku,
			TotalCount: 80,
			Reserved:   0,
		},
	}
)

type LomsServiceWithMock struct {
	handler             *LomsService
	orderRepositoryMock *mock.OrderRepositoryMock
	stockRepositoryMock *mock.StockRepositoryMock
}

func NewLomsServiceWithInMemoryRepository() *LomsService {
	orderRepository := orderrepo.NewOrderInMemoryRepository(10, &counter)

	stockRepository := stockrepo.NewStockInMemoryRepository(10)

	return NewLomsService(orderRepository, stockRepository)
}

func NewLomsServiceWithMock(t *testing.T) *LomsServiceWithMock {
	ctrl := minimock.NewController(t)

	orderRepository := mock.NewOrderRepositoryMock(ctrl)

	stockRepository := mock.NewStockRepositoryMock(ctrl)

	lomsService := NewLomsService(orderRepository, stockRepository)

	return &LomsServiceWithMock{
		handler:             lomsService,
		orderRepositoryMock: orderRepository,
		stockRepositoryMock: stockRepository,
	}
}

func TestHandler_All(t *testing.T) {
	t.Run("Test_OrderCancel", Test_OrderCancel)
	t.Run("Test_OrderCreate", Test_OrderCreate)
	t.Run("Test_OrderInfo", Test_OrderInfo)
	t.Run("Test_OrderPay", Test_OrderPay)
	t.Run("Test_StockInfo", Test_StockInfo)
}
