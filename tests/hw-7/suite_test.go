package hw7

import (
	"context"
	"testing"
	"time"

	"route256/tests/app/clients/loms"
	"route256/tests/app/domain"
	"route256/tests/app/kafka"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type lomsClient interface {
	OrderCreate(ctx context.Context, t provider.StepCtx, userID int64, items []domain.OrderItem) (orderID int64, statusCode int)
	OrderInfo(ctx context.Context, t provider.StepCtx, orderID int64) (order *domain.Order, statusCode int)
	OrderPay(ctx context.Context, t provider.StepCtx, orderID int64) (statusCode int)
	OrderCancel(ctx context.Context, t provider.StepCtx, orderID int64) (statusCode int)
	StocksInfo(ctx context.Context, t provider.StepCtx, sku int64) (count uint64, statusCode int)
}

type Suite struct {
	suite.Suite

	lomsClient    lomsClient
	kafkaConsumer *kafka.Consumer
	ctx           context.Context
	cancel        context.CancelFunc
}

func (s *Suite) BeforeAll(t provider.T) {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	//cfg, err := config.NewConfig()
	//t.Require().NoError(err, "Не удалось создать конфиг")

	s.lomsClient = loms.NewClient("http://localhost:8084")

	t.Log("Start consumer")
	s.kafkaConsumer = kafka.NewConsumer(
		s.ctx,
		[]string{"localhost:9092"},
		"loms.order-events",
		10*time.Second,
	)

	// Пропускаем все существующие сообщения в топике, т.к они не относятся к текущему тесту
	// Наверное это и не нужно (стоит удалить):
	// - у нас каждый раз в пайплайне стартует новая кафка с новым топиком (это может иметь смысл при локальном запуске)
	// - у нас установлен config.Consumer.Offsets.Initial = sarama.OffsetNewest
	t.Log("Skip messages")
	n := s.kafkaConsumer.ReadAllUntil(s.ctx, 10*time.Second)
	t.Logf("kafkaConsumer.ReadAllUntil skipped %d messages", len(n))
}

func (s *Suite) AfterAll(t provider.T) {
	defer s.cancel()

	if err := s.kafkaConsumer.Close(); err != nil {
		t.Log(err)
	}
}

func TestHW7Suite(t *testing.T) {
	suite.RunNamedSuite(t, "Домашнее задание 7", new(Suite))
}
