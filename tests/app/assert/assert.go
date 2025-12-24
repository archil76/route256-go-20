package assert

import (
	"route256/tests/app/domain"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func StatusCode(sCtx provider.StepCtx, expected, actual int) {
	sCtx.Require().Equal(expected, actual, "Не совпадает статус код")
}

func Cart(sCtx provider.StepCtx, expected, actual *domain.Cart) {
	sCtx.Require().NotNil(expected, "Ожидаемая корзина не должна быть nil")
	sCtx.Require().NotNil(actual, "Фактическая корзина не должна быть nil")

	// Есть какие-то проблемы при проверке uint типов через Require().Equal - выводятся адреса вместо значений при падении.
	sCtx.Require().Equal(int64(expected.TotalPrice), int64(actual.TotalPrice), "Не совпадает общая цена корзины")
	sCtx.Require().Equal(len(expected.Items), len(actual.Items), "Не совпадает количество товаров в корзине")

	for i := range len(expected.Items) {
		sCtx.Require().Equal(int64(expected.Items[i].SKU), int64(actual.Items[i].SKU), "Не совпадает SKU товара в корзине")
		sCtx.Require().Equal(expected.Items[i].Name, actual.Items[i].Name, "Не совпадает название товара в корзине")
		sCtx.Require().Equal(expected.Items[i].Count, actual.Items[i].Count, "Не совпадает количество товара в корзине")
		sCtx.Require().Equal(expected.Items[i].Price, actual.Items[i].Price, "Не совпадает цена товара в корзине")
	}
}

func EmptyCart(sCtx provider.StepCtx, cart *domain.Cart) {
	sCtx.Require().NotNil(cart, "Корзина не должна быть nil")
	sCtx.Require().Empty(cart.Items, "Ожидается пустая корзина")
}

func NotEmptyCart(sCtx provider.StepCtx, cart *domain.Cart) {
	sCtx.Require().NotNil(cart, "Корзина не должна быть nil")
	sCtx.Require().NotEmpty(cart.Items, "Ожидается наличие товаров в корзине")
}

func OrderID(sCtx provider.StepCtx, actual int64) {
	sCtx.Require().Greater(actual, int64(0), "Ожидается положительный идентификатор заказа")
}

func Order(sCtx provider.StepCtx, expected, actual *domain.Order) {
	sCtx.Require().Equal(expected.Status, actual.Status, "Не совпадает статус заказа")
	sCtx.Require().Equal(expected.User, actual.User, "Не совпадает пользователь заказа")
	sCtx.Require().Equal(len(expected.Items), len(actual.Items), "Не совпадает количество товаров в заказе")

	for i := range len(expected.Items) {
		sCtx.Require().Equal(int64(expected.Items[i].Sku), int64(actual.Items[i].Sku), "Не совпадает SKU товара в заказа")
		sCtx.Require().Equal(int64(expected.Items[i].Count), int64(actual.Items[i].Count), "Не совпадает количество товара в заказе")
	}
}

func OrderStatus(sCtx provider.StepCtx, expected, actual domain.OrderStatus) {
	sCtx.Require().Equal(expected, actual, "Не совпадает статус заказа")
}

func Stocks(sCtx provider.StepCtx, expected, actual uint64) {
	sCtx.Require().Equal(expected, actual, "Не совпадает количество стоков")
}

func OrderEvent(sCtx provider.StepCtx, expected, actual domain.OrderEvent) {
	sCtx.Require().Equal(expected.OrderID, actual.OrderID, "Не совпадает OrderID в событии")
	sCtx.Require().Equal(expected.Status, actual.Status, "Не совпадает статус заказа")
}
