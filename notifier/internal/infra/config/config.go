package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Kafka struct {
		Host            string `yaml:"host"`
		Port            string `yaml:"port"`
		OrderTopic      string `yaml:"order_topic"`       //nolint:revive
		ConsumerGroupID string `yaml:"consumer_group_id"` //nolint:revive
		Brokers         string `yaml:"brokers"`
	} `yaml:"kafka"`
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename) //nolint:gosec
	if err != nil {
		return nil, err
	}

	defer f.Close()

	config := &Config{}
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
