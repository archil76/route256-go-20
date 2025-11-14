package hw1

import (
	"context"
	"testing"

	"route256/tests/app/clients/cart"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type Suite struct {
	suite.Suite

	ctx context.Context

	cartClient *cart.Client
}

func TestSuite(t *testing.T) {
	suite.RunNamedSuite(t, "Домашнее задание 1", new(Suite))
}

func (s *Suite) BeforeAll(t provider.T) {

	s.ctx = context.Background()

	s.cartClient = cart.NewClient("http://localhost:8080")

}
