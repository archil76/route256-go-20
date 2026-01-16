package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/loms/internal/infra/logger"
	"strconv"
	"time"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

type KafkaMessage struct {
	OrderId int       `json:"order_id"`
	Status  string    `json:"status"`
	Moment  time.Time `json:"moment"`
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

func newKafkaMessage(orderID int64, status string) *KafkaMessage {
	kafkaMessage := KafkaMessage{}

	orderIDint := int(orderID) //nolint:gosec
	kafkaMessage.OrderId = orderIDint
	kafkaMessage.Status = status
	kafkaMessage.Moment = time.Now()

	return &kafkaMessage
}

func (p *KafkaProducer) SendMessage(orderID int64, status string) {
	kafkaMessage := newKafkaMessage(orderID, status)
	kafkaKey := strconv.FormatInt(orderID, 10)

	message, err := json.Marshal(kafkaMessage)
	if err != nil {
		logger.Errorw("Ошибка создания сообщения kafka", "топик", p.topic, "Error", err)
	}
	p.sendKafkaMessage(kafkaKey, message)
}

func (p *KafkaProducer) sendKafkaMessage(key string, message []byte) {
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
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
