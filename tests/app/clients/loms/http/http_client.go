package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type Client struct {
	baseUrl string
	cl      *http.Client
}

func NewClient(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
		cl:      &http.Client{},
	}
}

func (c *Client) OrderCreate(ctx context.Context, t provider.StepCtx, req *OrderCreateRequest) (*OrderCreateResponse, error) {
	data, err := json.Marshal(req)
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, data)

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/order/create", c.baseUrl), bytes.NewBuffer(data))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считываение ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	orderCreated := &OrderCreateResponse{}
	err = json.Unmarshal(body, &orderCreated)
	t.Require().NoError(err, "парсинг ответа")

	return orderCreated, nil
}

func (c *Client) OrderInfo(ctx context.Context, t provider.StepCtx, req *OrderInfoRequest) (*OrderInfoResponse, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/order/info?orderId=%d", c.baseUrl, req.OrderID), nil)
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "считывание ответа")
	t.WithNewAttachment("response body", allure.JSON, body)

	orderInfo := &OrderInfoResponse{}
	err = json.Unmarshal(body, &orderInfo)
	t.Require().NoError(err, "парсинг ответа")

	return orderInfo, nil
}

func (c *Client) OrderPay(ctx context.Context, t provider.StepCtx, req *OrderPayRequest) (*OrderPayResponse, error) {
	data, err := json.Marshal(req)
	t.Require().NoError(err, "не удалось сериализовать запрос")
	t.WithNewAttachment("request payload", allure.JSON, data)

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/order/pay", c.baseUrl), bytes.NewBuffer(data))
	t.Require().NoError(err, "не удалось создать запрос")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "не удалось выполнить запрос")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	orderPaid := &OrderPayResponse{}
	err = json.Unmarshal(body, &orderPaid)
	t.Require().NoError(err, "не удалось распарсить ответ")

	return orderPaid, nil
}

func (c *Client) OrderCancel(ctx context.Context, t provider.StepCtx, req *OrderCancelRequest) (*OrderCancelResponse, error) {
	data, err := json.Marshal(req)
	t.Require().NoError(err, "не удалось сериализовать запрос")
	t.WithNewAttachment("request payload", allure.JSON, data)

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/order/cancel", c.baseUrl), bytes.NewBuffer(data))
	t.Require().NoError(err, "не удалось создать запрос")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "не удалось выполнить запрос")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	orderCanceled := &OrderCancelResponse{}
	err = json.Unmarshal(body, &orderCanceled)
	t.Require().NoError(err, "не удалось распарсить ответ")

	return orderCanceled, nil
}

func (c *Client) StocksInfo(ctx context.Context, t provider.StepCtx, req *StocksInfoRequest) (*StocksInfoResponse, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/stock/info?sku=%d", c.baseUrl, req.Sku), nil)
	t.Require().NoError(err, "не удалось создать запрос")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "не удалось выполнить запрос")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	infoResponse := &StocksInfoResponse{}
	err = json.Unmarshal(body, &infoResponse)
	t.Require().NoError(err, "не удалось распарсить ответ")

	return infoResponse, nil
}
