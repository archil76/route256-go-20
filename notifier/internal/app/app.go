package app

import (
	"context"
	"fmt"
	"route256/notifier/internal/infra/config"
	kafkaConsumer "route256/notifier/internal/infra/kafka/consumer"
	"route256/notifier/internal/infra/logger"
)

type App struct {
	config        *config.Config
	consumerGroup *kafkaConsumer.KafkaConsumerGroup
}

func NewApp(configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	ctx := context.Background()

	consumerGroup, err := kafkaConsumer.NewConsumerGroup(ctx, c.Kafka.Brokers, c.Kafka.OrderTopic, c.Kafka.ConsumerGroupID)
	if err != nil {
		return nil, err
	}

	app := &App{config: c, consumerGroup: &consumerGroup}

	return app, nil
}

func (app *App) ListenAndServe() error {
	defer app.consumerGroup.Close()

	//go func() {
	//	// слушаем ошибки группы в отдельной горутине
	//	for err := range app.consumerGroup.Errors() {
	//		logger.Errorw("ConsumerGroup error:", "err", err)
	//	}
	//}()
	fmt.Printf("notifier service is ready %s:%s\n", app.config.Kafka.Host, app.config.Kafka.Port)

	for {
		ctx := context.Background()

		// Метод Consume запускает обработку; после ребалансировки может возвращаться и вызываться снова
		err := app.consumerGroup.Consume(ctx)
		if err != nil {
			logger.Fatalw("ConsumerGroup error:", "err", err, "topic", app.config.Kafka.OrderTopic, "consumerGroup", app.config.Kafka.ConsumerGroupID)
		}

		// Если пришел сигнал прекращения, выходим из цикла
		if ctx.Err() != nil {
			break
		}
	}

	return nil
}
