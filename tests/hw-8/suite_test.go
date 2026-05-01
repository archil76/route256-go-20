package hw_8

import (
	"context"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"route256/tests/app/clients/comments"
	//"route256/tests/app/config"
)

const (
	defaultEditInterval = time.Second
)

type Suite struct {
	suite.Suite

	ctx context.Context

	commentsClient *comments.Client
	editInterval   time.Duration
}

func TestSuite(t *testing.T) {
	suite.RunNamedSuite(t, "Домашнее задание 8", new(Suite))
}

func (s *Suite) BeforeAll(t provider.T) {
	//cfg, err := config.NewConfig()
	//t.Require().NoError(err, "Не удалось создать конфиг")

	s.ctx = context.Background()

	editInterval := defaultEditInterval
	if cfgEditInterval, err := time.ParseDuration("1s"); err != nil {
		editInterval = cfgEditInterval
	}
	t.WithNewParameters("editInterval", editInterval.String())

	s.editInterval = editInterval
	//s.commentsClient = comments.NewClient(cfg.Env.CommentsServiceUrl)
	s.commentsClient = comments.NewClient("http://localhost:8086")
}
