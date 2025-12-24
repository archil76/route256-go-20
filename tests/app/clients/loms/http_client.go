package loms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

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

func (c *Client) OrderCreate(
	ctx context.Context,
	t provider.StepCtx,
	userID int64,
	items []domain.OrderItem,
) (orderID int64, statusCode int) {
	req := orderCreateRequest{
		UserID: userID,
		Items:  make([]orderItem, 0, len(items)),
	}
	for _, item := range items {
		req.Items = append(req.Items, orderItem{
			Sku:   fmt.Sprintf("%d", item.Sku),
			Count: item.Count,
		})
	}
	data, err := json.Marshal(req)
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request body", allure.JSON, data)

	url := fmt.Sprintf("%s/order/create", c.baseUrl)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	orderResp := &orderCreateResponse{}
	err = json.Unmarshal(body, orderResp)
	t.Require().NoError(err, "парсинг ответа")

	return orderResp.OrderID, res.StatusCode
}

func (c *Client) OrderInfo(ctx context.Context, t provider.StepCtx, orderID int64) (order *domain.Order, statusCode int) {
	url := fmt.Sprintf("%s/order/info?orderId=%d", c.baseUrl, orderID)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	orderInfoResp := &OrderInfoResponse{}
	err = json.Unmarshal(body, orderInfoResp)
	t.Require().NoError(err, "парсинг ответа")

	order = &domain.Order{
		Status: domain.OrderStatus(orderInfoResp.Status),
		User:   orderInfoResp.User,
		Items:  make([]domain.OrderItem, 0, len(orderInfoResp.Items)),
	}
	for _, item := range orderInfoResp.Items {
		sku, err := strconv.ParseInt(item.Sku, 10, 64)
		t.Require().NoError(err, "parse sku")
		order.Items = append(order.Items, domain.OrderItem{
			Sku:   int32(sku),
			Count: item.Count,
		})
	}

	return order, res.StatusCode
}

func (c *Client) OrderPay(ctx context.Context, t provider.StepCtx, orderID int64) (statusCode int) {
	data, err := json.Marshal(orderRequest{OrderID: orderID})
	t.Require().NoError(err, "не удалось сериализовать запрос")
	t.WithNewAttachment("request payload", allure.JSON, data)

	url := fmt.Sprintf("%s/order/pay", c.baseUrl)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	t.Require().NoError(err, "не удалось создать запрос")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "не удалось выполнить запрос")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	return res.StatusCode
}

func (c *Client) OrderCancel(ctx context.Context, t provider.StepCtx, orderID int64) (statusCode int) {
	data, err := json.Marshal(orderRequest{OrderID: orderID})
	t.Require().NoError(err, "не удалось сериализовать запрос")
	t.WithNewAttachment("request payload", allure.JSON, data)

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/order/cancel", c.baseUrl), bytes.NewBuffer(data))
	t.Require().NoError(err, "не удалось создать запрос")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "не удалось выполнить запрос")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	return res.StatusCode
}

func (c *Client) StocksInfo(ctx context.Context, t provider.StepCtx, sku int64) (count uint64, statusCode int) {
	url := fmt.Sprintf("%s/stock/info?sku=%d", c.baseUrl, sku)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	t.Require().NoError(err, "не удалось создать запрос")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "не удалось выполнить запрос")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	stocksInfoResp := &stocksInfoResponse{}
	err = json.Unmarshal(body, &stocksInfoResp)
	t.Require().NoError(err, "не удалось распарсить ответ")

	return stocksInfoResp.Count, res.StatusCode
}
