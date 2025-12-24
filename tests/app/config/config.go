package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// EnvConfig конфиг из переменных окружения
type EnvConfig struct {
	CartServiceUrl               string
	LomsServiceUrl               string
	KafkaBrokers                 []string
	KafkaTopic                   string
	KafkaWaitingForEventDuration time.Duration
	KafkaBeforeAllSkipDuration   time.Duration
	CommentsServiceUrl           string
	CommentsServiceEditInterval  string
	ProductsServiceUrl           string
}

type Config struct {
	Env *EnvConfig
}

func NewConfig() (*Config, error) {
	//err := godotenv.Load("../.env")
	//if err != nil {
	//	return nil, err
	//}

	kafkaWaitingForEventDuration, err := time.ParseDuration(getEnv("KAFKA_WAITING_FOR_EVENT_DURATION", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid KAFKA_WAITING_FOR_EVENT_DURATION env, unable to parse duration: %w", err)
	}
	kafkaBeforeAllSkipDuration, err := time.ParseDuration(getEnv("KAFKA_BEFORE_ALL_SKIP_DURATION", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid KAFKA_BEFORE_ALL_SKIP_DURATION env, unable to parse duration: %w", err)
	}

	return &Config{
		Env: &EnvConfig{
			CartServiceUrl:               getEnv("CART_SERVICE_URL", ""),
			LomsServiceUrl:               getEnv("LOMS_SERVICE_URL", ""),
			KafkaBrokers:                 getEnvs("KAFKA_BROKERS", []string{}),
			KafkaTopic:                   getEnv("KAFKA_TOPIC", ""),
			KafkaWaitingForEventDuration: kafkaWaitingForEventDuration,
			KafkaBeforeAllSkipDuration:   kafkaBeforeAllSkipDuration,
			CommentsServiceUrl:           getEnv("COMMENTS_SERVICE_URL", ""),
			CommentsServiceEditInterval:  getEnv("COMMENTS_SERVICE_EDIT_INTERVAL", ""),
			//ProductsServiceUrl:           getEnv("PRODUCT_SERVICE_URL", ""),
		},
	}, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvs(key string, defaultVals []string) []string {
	if values, exists := os.LookupEnv(key); exists {
		return strings.Split(values, ",")
	}

	return defaultVals
}
