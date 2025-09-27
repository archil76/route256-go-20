package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host     string `yaml:"host"`
		GrpcPort string `yaml:"grpc_port"`
		HttpPort string `yaml:"http_port"` //nolint:revive
	} `yaml:"service"`

	Jaeger struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"jaeger"`

	DBMaster struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
	} `yaml:"db_master"`

	DBReplica struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
	} `yaml:"db_replica"`

	Kafka struct {
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		OrderTopic string `yaml:"order_topic"`
		Brokers    string `yaml:"brokers"`
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
