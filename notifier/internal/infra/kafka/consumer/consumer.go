package consumer

import (
	"route256/notifier/internal/infra/logger"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct{}

func (c *KafkaConsumer) Setup(sess sarama.ConsumerGroupSession) error {
	logger.Infow("Участник группы запущен", "MemberID", sess.MemberID())
	return nil
}
func (c *KafkaConsumer) Cleanup(sess sarama.ConsumerGroupSession) error {
	logger.Infow("Участник группы завершает работу", "MemberID", sess.MemberID())
	return nil
}
func (c *KafkaConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger.Infow("Чтение партиции", "Partition", claim.Partition())

	for msg := range claim.Messages() {
		logger.Infow("Получено сообщение", "Topic", msg.Topic, "Offset", msg.Offset, "Key", msg.Key, "Value", string(msg.Value))
		// помечаем сообщение как прочитанное
		sess.MarkMessage(msg, "")
	}

	return nil // выходим, когда внешний контекст завершён
}
