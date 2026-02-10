package outbox

import (
	mock2 "route256/loms/internal/domain/service/outbox/mock"
	"testing"

	"github.com/gojuno/minimock/v3"
	"golang.org/x/net/context"
)

type OutboxServiceWithMock struct {
	handler              *OutboxService
	outboxRepositoryMock *mock2.OutboxRepositoryMock
	producerMock         *mock2.KafkaProducerMock
	poolerMock           *mock2.PgPoolerMock
}

func NewOutboxServiceWithMock(t *testing.T) *OutboxServiceWithMock {
	ctrl := minimock.NewController(t)

	outboxRepository := mock2.NewOutboxRepositoryMock(ctrl)
	producer := mock2.NewKafkaProducerMock(ctrl)
	pooler := mock2.NewPgPoolerMock(ctrl)

	outboxService := NewOutboxService(context.Background(), outboxRepository, 1, producer, pooler)

	return &OutboxServiceWithMock{
		handler:              outboxService,
		outboxRepositoryMock: outboxRepository,
		producerMock:         producer,
		poolerMock:           pooler,
	}
}

func TestHandler_All(t *testing.T) {
	t.Run("Test_CreateMessage", Test_CreateMessage)
	t.Run("Test_SendMessages", Test_SendMessages)

}
