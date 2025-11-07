package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Workers string `yaml:"workers"`
	} `yaml:"service"`

	Jaeger struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"jaeger"`

	ProductService struct {
		Host  string `yaml:"host"`
		Port  string `yaml:"port"`
		Token string `yaml:"token"`
		Limit string `yaml:"limit"`
		Burst string `yaml:"burst"`
	} `yaml:"product_service"`

	LomsService struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"loms_service"`
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
