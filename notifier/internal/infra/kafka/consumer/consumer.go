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

type KafkaMessage struct {
	OrderId              int64  `json:"order_id"`
	Status               string `json:"status"`
	Moment               string `json:"moment"`
	AdditionalProperties bool   `json:"additionalProperties"`
}

func NewConsumer(_ context.Context, addrs, topic, groupID string) (KafkaConsumerGroup, error) {
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
	return cg.Errors()
}

type KafkaConsumer struct{}

func (c *KafkaConsumer) Setup(sess sarama.ConsumerGroupSession) error {
	logger.Infow("[Consumer] Участник группы запущен, ассайн получен", "MemberID", sess.MemberID())
	return nil
}
func (c *KafkaConsumer) Cleanup(sess sarama.ConsumerGroupSession) error {
	logger.Infow("[Consumer] Участник группы завершает работу", "MemberID", sess.MemberID())
	return nil
}
func (c *KafkaConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger.Infow("-> [Consumer] Чтение партиции", "Partition", claim.Partition())

	for msg := range claim.Messages() {
		logger.Infow(fmt.Sprintf("[Consumer] %s: раздел=%d офсет=%d ключ=%s значение=%s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value)))
		// помечаем сообщение как прочитанное
		sess.MarkMessage(msg, "")
	}

	return nil // выходим, когда внешний контекст завершён
}
