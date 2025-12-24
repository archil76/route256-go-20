package hw3

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"route256/tests/app/clients/cart"
	"route256/tests/app/clients/loms"
)

// Перед запуском тестов следует сверить данные по стокам в бд
// с данными из docs/homework-3/stock-data.json.

type Suite struct {
	suite.Suite

	cartClient *cart.Client
	lomsClient *loms.Client
}

func TestSuite(t *testing.T) {
	suite.RunNamedSuite(t, "Домашнее задание 3", new(Suite))
}

func (s *Suite) BeforeAll(t provider.T) {
	//cfg, err := config.NewConfig()
	//t.Require().NoError(err, "Не удалось создать конфиг")

	s.cartClient = cart.NewClient("http://localhost:8080")
	s.lomsClient = loms.NewClient("http://localhost:8084")
}
