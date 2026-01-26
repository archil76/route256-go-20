package outbox

import (
	"context"
	"route256/loms/internal/domain/model"
	"time"
)

type OutboxRepository interface {
	Create(ctx context.Context, key, status string, payload []byte) (int, error)
	SetStatus(ctx context.Context, id int, status string) error
	Get(ctx context.Context) (*[]model.OutboxItem, error)
}

type KafkaProducer interface {
	SendMessage(ctx context.Context, key string, message []byte) error
}

type Status string

const (
	NEWSTATUS Status = "new"
	SENT      Status = "sent"
	ERROR     Status = "error"
)

type KafkaMessage struct {
	OrderId int       `json:"order_id"`
	Status  string    `json:"status"`
	Moment  time.Time `json:"moment"`
}

type OutboxService struct {
	outboxRepository OutboxRepository
	producer         KafkaProducer
}

func NewOutboxService(ctx context.Context, outboxRepository OutboxRepository, interval int, producer KafkaProducer) *OutboxService {
	s := OutboxService{
		outboxRepository: outboxRepository,
		producer:         producer,
	}

	go func() {
		t := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			select {
			case <-t.C:
				s.SendMessages(ctx)
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}()

	return &s
}

func newKafkaMessage(orderID int64, status string) *KafkaMessage {
	kafkaMessage := KafkaMessage{}

	orderIDint := int(orderID) //nolint:gosec
	kafkaMessage.OrderId = orderIDint
	kafkaMessage.Status = status
	kafkaMessage.Moment = time.Now()

	return &kafkaMessage
}
