package kafka

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

const (
	groupID = "integration-tests-group"
)

type Consumer struct {
	ch                      <-chan *sarama.ConsumerMessage
	cg                      sarama.ConsumerGroup
	waitingForEventDuration time.Duration
}

func (c *Consumer) Close() error {
	return errors.Wrap(c.cg.Close(), "ConsumerGroup.Close")
}

func (c *Consumer) ConsumeSingle(ctx context.Context) *sarama.ConsumerMessage {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.waitingForEventDuration)
	defer cancel()

	return c.consumeSingle(timeoutCtx)
}

func (c *Consumer) consumeSingle(ctx context.Context) *sarama.ConsumerMessage {
	select {
	case <-ctx.Done():
		break
	case msg, ok := <-c.ch:
		if ctx.Err() != nil {
			break
		}
		if !ok {
			break
		}
		return msg
	}

	return nil
}

func (c *Consumer) ReadAllUntil(ctx context.Context, d time.Duration) []*sarama.ConsumerMessage {
	timeoutCtx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	var msgs []*sarama.ConsumerMessage
	for {
		msg := c.ConsumeSingle(timeoutCtx)
		if msg != nil {
			msgs = append(msgs, msg)
		} else {
			break
		}
	}

	return msgs
}

func NewConsumer(ctx context.Context, brokers []string, topic string, waitingForEventDuration time.Duration) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	// наша консьюмер-группа одноразовая и не подразумевает переподключения
	config.Consumer.Offsets.AutoCommit.Enable = false

	// подмешиваем к cg случайную соль на случай если будет переиспользован контейнер с кафкой
	group, err := sarama.NewConsumerGroup(brokers, fmt.Sprintf("%s-%d", groupID, rand.Int63()), config)
	if err != nil {
		panic(err)
	}

	ch := make(chan *sarama.ConsumerMessage)
	handler := consumerGroupHandler{
		ch: ch,
	}

	c := &Consumer{
		ch:                      ch,
		cg:                      group,
		waitingForEventDuration: waitingForEventDuration,
	}

	go func() {
		for grErr := range group.Errors() {
			fmt.Println(grErr)
		}
	}()

	go func() {
		if consumeErr := group.Consume(ctx, []string{topic}, handler); consumeErr != nil {
			fmt.Println(consumeErr)
		}
	}()

	return c
}

type consumerGroupHandler struct {
	ch chan *sarama.ConsumerMessage
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			fmt.Printf("Message topic:%q partition:%d offset:%d\n", message.Topic, message.Partition, message.Offset)
			h.ch <- message

			// mark message as successfully handled and ready to commit offset
			// autocommit may commit message offset sometime
			sess.MarkMessage(message, "")
		case <-sess.Context().Done():
			return nil
		}
	}

	return nil
}
