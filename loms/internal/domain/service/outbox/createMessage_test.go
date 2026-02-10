package outbox

import (
	"context"
	"errors"
	"route256/loms/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateMessage(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()

	testHandler := NewOutboxServiceWithMock(t)

	t.Run("Добавление сообщения. Успешный путь", func(t *testing.T) {
		var err error
		var orderID int64

		handler := testHandler.handler
		orderID = int64(10)

		testHandler.outboxRepositoryMock.CreateMock.Set(func(ctx context.Context, key string, status string, payload []byte) (id int, err error) {

			assert.Equal(t, key, "10")

			if status != "new" && status != "awaiting payment" {
				return 0, errors.New("Ожидается статус new или awaiting payment")
			}
			assert.IsType(t, []byte{}, payload)

			return 1, nil
		})

		testHandler.poolerMock.InTxMock.Set(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})

		handler.CreateMessage(ctx, orderID, model.NEWSTATUS)

		require.NoError(t, err)

		handler.CreateMessage(ctx, orderID, model.AWAITINGPAYMENT)

		require.NoError(t, err)
	})
}
