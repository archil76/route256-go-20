package producer

import (
	"context"
	"fmt"
	"route256/loms/internal/infra/logger"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(_ context.Context, addrs, topic string) (KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	producer, err := sarama.NewSyncProducer([]string{addrs}, config)
	if err != nil {
		return KafkaProducer{}, fmt.Errorf("не удалось подключиться к Kafka: %w", err)
	}

	return KafkaProducer{producer, topic}, nil
}

func (p *KafkaProducer) SendMessage(_ context.Context, key string, message []byte) error {
	return p.sendKafkaMessage(key, message)
}

func (p *KafkaProducer) sendKafkaMessage(key string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Partition: 0,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.StringEncoder(message),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		logger.Errorw("Ошибка при отправке сообщения в топик", "топик", p.topic, "Error", err)
	}

	return err
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
