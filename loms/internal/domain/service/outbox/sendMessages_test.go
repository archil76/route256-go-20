package outbox

import (
	"context"
	"route256/loms/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SendMessages(t *testing.T) {
	t.Parallel()
	t.Helper()

	ctx := context.Background()

	testHandler := NewOutboxServiceWithMock(t)

	t.Run("Отправка сообщений. Успешный путь", func(t *testing.T) {
		var err error
		var key string
		var message []byte

		handler := testHandler.handler
		key = "10"
		message = []byte("message")
		item := model.OutboxItem{
			Id:      1,
			Key:     key,
			Payload: message,
		}
		items := []model.OutboxItem{item}

		testHandler.outboxRepositoryMock.GetMock.When(ctx).Then(&items, nil)
		testHandler.producerMock.SendMessageMock.When(ctx, items[0].Key, items[0].Payload).Then(nil)
		testHandler.outboxRepositoryMock.SetStatusMock.When(ctx, int(items[0].Id), "sent").Then(nil)
		testHandler.poolerMock.InTxMock.Set(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})

		handler.SendMessages(ctx)

		require.NoError(t, err)

	})
}
