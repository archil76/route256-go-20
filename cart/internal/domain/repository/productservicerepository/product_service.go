package productservicerepository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"route256/cart/internal/domain/model"
	"route256/cart/internal/infra/errgroup"
	"route256/cart/internal/infra/ratelimiter"
	"time"
)

var (
	ErrNotOk           = errors.New("status not ok")
	ErrProductNotFound = errors.New("product not found")
)

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
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/product/%d", s.address, sku),
		http.NoBody,
	)
	if err != nil {
		return nil, ErrNotOk
	}

	req.Header.Add("X-API-KEY", s.token)

	fmt.Printf("Request to ProductService: %s %d\n", s.address, sku)

	response, err := s.httpClient.Do(req)
	if err != nil {
		return nil, ErrNotOk
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrProductNotFound
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrNotOk
	}

	resp := &GetProductResponse{}
	if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
		return nil, ErrNotOk
	}

	return &model.Product{
		Name:  resp.Name,
		Price: uint32(resp.Price),
		Sku:   resp.Sku,
	}, nil
}

func (s *ProductService) GetProductsBySkus(ctx context.Context, skus []model.Sku) ([]model.Product, error) {
	fmt.Printf("GetProductsBySkus start\n")

	limit := 10
	duration := time.Second

	group, ctx := errgroup.WithContext(ctx)
	rateLimiter, err := ratelimiter.WithContext(ctx, limit, duration)
	if err != nil {
		return nil, err
	}
	rateLimiter.Start()
	defer rateLimiter.Stop()

	products := make([]model.Product, len(skus))
	for i, sku := range skus {
		group.Go(func() error {
			rateLimiter.Wait()
			product, err := s.GetProductBySku(ctx, sku)
			if err != nil {
				ctx.Done()
				return err
			}

			products[i] = *product

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	fmt.Printf("GetProductsBySkus stop\n")
	return products, nil
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price int32  `json:"price"`
	Sku   int64  `json:"sku"`
}
