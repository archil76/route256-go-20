package cart

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/domain"
)

type Client struct {
	baseUrl string
	cl      *http.Client
}

func NewClient(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
		cl:      http.DefaultClient,
	}
}

func (c *Client) AddItem(ctx context.Context, t provider.StepCtx, userID, sku int64, count int64) int {
	data, err := json.Marshal(addItemRequest{
		Count: count,
	})
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, data)

	url := fmt.Sprintf("%s/user/%d/cart/%d", c.baseUrl, userID, sku)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	return res.StatusCode
}

func (c *Client) DeleteItem(ctx context.Context, t provider.StepCtx, userID, sku int64) int {
	url := fmt.Sprintf("%s/user/%d/cart/%d", c.baseUrl, userID, sku)
	r, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	t.Require().NoError(err, "создание запроса")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	return res.StatusCode
}

func (c *Client) GetCart(ctx context.Context, t provider.StepCtx, userID int64) (*domain.Cart, int) {
	url := fmt.Sprintf("%s/user/%d/cart", c.baseUrl, userID)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	t.Require().NoError(err, "создание запроса")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	if res.StatusCode != http.StatusOK {
		return &domain.Cart{}, res.StatusCode
	}

	listResp := &listResponse{}
	err = json.Unmarshal(body, listResp)
	t.Require().NoError(err, "парсинг ответа")

	cart := &domain.Cart{
		Items:      make([]domain.CartItem, 0, len(listResp.Items)),
		TotalPrice: listResp.TotalPrice,
	}
	for _, item := range listResp.Items {
		cart.Items = append(cart.Items, domain.CartItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return cart, res.StatusCode
}

func (c *Client) Checkout(ctx context.Context, t provider.StepCtx, userID int64) (orderID int64, statusCode int) {
	data, err := json.Marshal(checkoutRequest{
		User: userID,
	})
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, data)

	url := fmt.Sprintf("%s/checkout/%d", c.baseUrl, userID)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, res.StatusCode
	}

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	checkoutResp := &checkoutResponse{}
	err = json.Unmarshal(body, checkoutResp)
	t.Require().NoError(err, "парсинг ответа")

	return checkoutResp.OrderID, res.StatusCode
}

func (c *Client) DeleteCart(ctx context.Context, t provider.StepCtx, userID int64) (statusCode int) {
	url := fmt.Sprintf("%s/user/%d/cart", c.baseUrl, userID)
	r, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	t.Require().NoError(err, "создание запроса")

	t.WithNewParameters("request", fmt.Sprintf("%s %s", r.Method, r.URL))

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	return res.StatusCode
}
