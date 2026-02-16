package loms

import (
	"context"
	"route256/loms/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_OrderCreate(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()

	testHandler := NewLomsServiceWithMock(t)

	t.Run("Добавление Заказа. Успешный путь", func(t *testing.T) {
		var err error
		var orderIDRes int64

		testHandler.orderRepositoryMock.CreateMock.Expect(ctx, newOrder).Return(&order, nil)
		testHandler.outboxServiceMock.CreateMessageMock.Set(func(ctx context.Context, orderID int64, status model.Status) {})
		testHandler.poolerMock.InTxMock.Set(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})
		testHandler.stockRepositoryMock.ReserveMock.Expect(ctx, items).Return(stocks, nil)
		testHandler.orderRepositoryMock.SetStatusMock.Expect(ctx, awaitingPaymentOrder, model.AWAITINGPAYMENT).Return(nil)

		handler := testHandler.handler

		orderIDRes, err = handler.OrderCreate(ctx, tp.userID, items)

		require.NoError(t, err)
		require.Equal(t, tp.orderID, orderIDRes)

		var orderID2Res int64

		testHandler.orderRepositoryMock.CreateMock.Expect(ctx, newOrder2).Return(&order2, nil)
		testHandler.outboxServiceMock.CreateMessageMock.Set(func(ctx context.Context, orderID2 int64, status model.Status) {})
		testHandler.poolerMock.InTxMock.Set(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})
		testHandler.stockRepositoryMock.ReserveMock.Expect(ctx, items).Return(stocks, nil)
		testHandler.orderRepositoryMock.SetStatusMock.Expect(ctx, awaitingPaymentOrder2, model.AWAITINGPAYMENT).Return(nil)

		orderID2Res, err = handler.OrderCreate(ctx, tp.userID2, items)

		require.NoError(t, err)
		require.Equal(t, tp.orderID2, orderID2Res)
	})
}

//func TestOrderCreate_SuccessAsync(t *testing.T) {
//	t.Parallel()
//	t.Helper()
//
//	const (
//		sku   = 1625903
//		count = 1
//
//		ordersCount = 100
//	)
//
//	var (
//		userIDs = []int64{42, 43, 44}
//
//		ctx = context.Background()
//
//		testHandler = NewLomsServiceWithMock(t)
//	)
//
//	t.Run("Успешное создание заказов (async)", func(t *testing.T) {
//
//		testHandler.orderRepositoryMock.CreateMock.Set(func(ctx context.Context, order model.Order) (op1 *model.Order, err error) {
//			return &order, nil
//		})
//		testHandler.outboxServiceMock.CreateMessageMock.Set(func(ctx context.Context, orderID int64, status model.Status) {})
//		testHandler.poolerMock.InTxMock.Set(func(ctx context.Context, fn func(ctx context.Context) error) error {
//			return fn(ctx)
//		})
//		testHandler.stockRepositoryMock.ReserveMock.Set(func(ctx context.Context, items []model.Item) (sa1 []model.Stock, err error) {
//			return nil, nil
//		})
//		testHandler.orderRepositoryMock.SetStatusMock.Set(func(ctx context.Context, order model.Order, status model.Status) (err error) {
//			return nil
//		})
//
//		var wg sync.WaitGroup
//
//		for i := range ordersCount {
//			wg.Add(1)
//			go func() {
//				defer wg.Done()
//
//				var err error
//				var orderIDRes int64
//
//				userID := userIDs[i%len(userIDs)]
//				handler := testHandler.handler
//
//				orderIDRes, err = handler.OrderCreate(ctx, userID, []model.Item{
//					{Sku: sku, Count: count},
//				})
//
//				require.NoError(t, err)
//				require.Greater(t, 0, orderIDRes)
//
//			}()
//		}
//		wg.Wait()
//	})
//}
