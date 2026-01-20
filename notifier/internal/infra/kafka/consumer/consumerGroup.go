package consumer

import (
	"context"
	"fmt"
	"route256/notifier/internal/infra/logger"
	"time"

	"github.com/IBM/sarama"
)

type KafkaConsumerGroup struct {
	consumerGroup sarama.ConsumerGroup
	topics        []string
}

func NewConsumerGroup(_ context.Context, addrs, topic, groupID string) (KafkaConsumerGroup, error) {
	config := sarama.NewConfig()

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second

	consumerGroup, err := sarama.NewConsumerGroup([]string{addrs}, groupID, config)
	if err != nil {
		return KafkaConsumerGroup{}, fmt.Errorf("не удалось подключиться к Kafka: %w", err)
	}

	topics := []string{topic}

	return KafkaConsumerGroup{consumerGroup, topics}, nil
}

func (cg *KafkaConsumerGroup) Consume(ctx context.Context) error {

	return cg.consumerGroup.Consume(ctx, cg.topics, &KafkaConsumer{})
}

func (cg *KafkaConsumerGroup) Close() {
	err := cg.consumerGroup.Close()
	if err != nil {
		logger.Infow("Ошибка закрытия consumer group:", "error", err)
		return
	}
}

func (cg *KafkaConsumerGroup) Errors() <-chan error {
	return cg.consumerGroup.Errors()
}
