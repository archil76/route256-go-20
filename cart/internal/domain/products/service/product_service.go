package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"route256/cart/internal/domain/model"
)

var ErrNotOk = errors.New("status not ok")

type ProductService struct {
	httpClient http.Client
	token      string
	address    string
}

func NewProductService(httpClient http.Client, token string, address string) *ProductService {
	return &ProductService{
		httpClient: httpClient,
		token:      token,
		address:    address,
	}
}

func (s *ProductService) GetProductBySku(ctx context.Context, sku model.Sku) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/product/%d", s.address, sku),
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Add("X-API-KEY", s.token)

	response, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpClient.Do: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, model.ErrProductNotFound
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrNotOk
	}

	resp := &GetProductResponse{}
	if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
		return nil, fmt.Errorf("json.NewDecoder: %w", err)
	}

	return &model.Product{
		Name:  resp.Name,
		Price: resp.Price,
		Sku:   model.Sku(resp.Sku),
	}, nil
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price int32  `json:"price"`
	Sku   int64  `json:"sku"`
}
