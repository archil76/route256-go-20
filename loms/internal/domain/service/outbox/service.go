package outbox

import (
	"context"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/infra/logger"
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

type PgPooler interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
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
	pooler           PgPooler
	ticker           time.Ticker
	ctx              context.Context
	cancel           context.CancelFunc
}

func NewOutboxService(ctx context.Context, outboxRepository OutboxRepository, interval int, producer KafkaProducer, pooler PgPooler) *OutboxService {
	ctx, cancel := context.WithCancel(ctx)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	s := OutboxService{
		outboxRepository: outboxRepository,
		producer:         producer,
		pooler:           pooler,
		ticker:           *ticker,
		ctx:              ctx,
		cancel:           cancel,
	}

	return &s
}

func (s *OutboxService) Start() {
	go func() {
		defer logger.Infow("Sender goroutine stopped")

		for {
			select {
			case <-s.ctx.Done():
				logger.Infow("Stopping sender")
				return
			case <-s.ticker.C:
				if s.ctx.Err() != nil {
					return
				}
				s.SendMessages(s.ctx)
			}
		}

	}()
}

func (s *OutboxService) Stop() {
	s.ticker.Stop()
	s.cancel()
}

func newKafkaMessage(orderID int64, status string) *KafkaMessage {
	kafkaMessage := KafkaMessage{}

	orderIDint := int(orderID) //nolint:gosec
	kafkaMessage.OrderId = orderIDint
	kafkaMessage.Status = status
	kafkaMessage.Moment = time.Now()

	return &kafkaMessage
}
