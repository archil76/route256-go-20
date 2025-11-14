package hw4

import (
	"context"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"route256/tests/app/clients/cart"
	"route256/tests/app/clients/loms"
	"route256/tests/app/domain"
)

// Перед запуском тестов следует сверить данные по стокам в бд (миграции)
// с данными из docs/homework-4/stock-data.json.

type lomsClient interface {
	OrderCreate(ctx context.Context, t provider.StepCtx, userID int64, items []domain.OrderItem) (orderID int64, statusCode int)
	OrderInfo(ctx context.Context, t provider.StepCtx, orderID int64) (order *domain.Order, statusCode int)
	OrderPay(ctx context.Context, t provider.StepCtx, orderID int64) (statusCode int)
	OrderCancel(ctx context.Context, t provider.StepCtx, orderID int64) (statusCode int)
	StocksInfo(ctx context.Context, t provider.StepCtx, sku int64) (count uint64, statusCode int)
}

type cartClient interface {
	AddItem(ctx context.Context, t provider.StepCtx, userID, sku int64, count int64) int
	DeleteItem(ctx context.Context, t provider.StepCtx, userID, sku int64) int
	GetCart(ctx context.Context, t provider.StepCtx, userID int64) (*domain.Cart, int)
	Checkout(ctx context.Context, t provider.StepCtx, userID int64) (orderID int64, statusCode int)
	DeleteCart(ctx context.Context, t provider.StepCtx, userID int64) (statusCode int)
}

type Suite struct {
	suite.Suite

	cartClient cartClient
	lomsClient lomsClient
}

func (s *Suite) BeforeAll(t provider.T) {
	//cfg, err := config.NewConfig()
	//t.Require().NoError(err, "Не удалось создать конфиг")
	//
	//s.cartClient = cart.NewClient(cfg.Env.CartServiceUrl)
	//s.lomsClient = loms.NewClient(cfg.Env.LomsServiceUrl)

	s.cartClient = cart.NewClient("http://localhost:8080")
	s.lomsClient = loms.NewClient("http://localhost:8084")
}

func TestHW4Suite(t *testing.T) {
	suite.RunNamedSuite(t, "Домашнее задание 4", new(Suite))
}
