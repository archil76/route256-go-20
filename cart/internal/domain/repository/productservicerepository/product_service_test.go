package productservicerepository

import (
	"context"
	"fmt"
	"net/http"
	"route256/cart/internal/domain/model"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHandler_All(t *testing.T) {
	t.Run("TestHandler_Product_Service_Product", TestHandler_Product_Service_Product)
	t.Run("TestHandler_Product_Service_Products", TestHandler_Product_Service_Products)
	t.Run("TestHandler_Product_Service_ProductsWithWrongSku", TestHandler_Product_Service_ProductsWithWrongSku)
}

func TestHandler_Product_Service_Product(t *testing.T) {
	ctx := context.Background()
	skus := []model.Sku{
		1076963,
		1148162,
		2618151,
		2956315,
		2958025,
		3596599,
		4465995,
		4288068,
	}
	transport := http.DefaultTransport

	httpClient := http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	productService := NewProductService(
		httpClient,
		"testToken",
		fmt.Sprintf("%s:%s", "http://localhost", "8082"),
		10,
	)

	for _, sku := range skus {
		product, err := productService.GetProductBySku(ctx, sku)
		require.NoError(t, err)
		require.NotNil(t, product)
	}

}

func TestHandler_Product_Service_Products(t *testing.T) {
	ctx := context.Background()
	skus := []model.Sku{
		1625903,
		1148162,
		1076963,
		32638658,
		32605854,
		32205848,
	}
	transport := http.DefaultTransport

	httpClient := http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	productService := NewProductService(
		httpClient,
		"testToken",
		fmt.Sprintf("%s:%s", "http://localhost", "8082"),
		10,
	)

	products, err := productService.GetProductsBySkus(ctx, skus)
	require.NoError(t, err)
	require.Equal(t, len(skus), len(products))
}

func TestHandler_Product_Service_ProductsWithWrongSku(t *testing.T) {
	ctx := context.Background()
	wrongSku := model.Sku(32205849)
	skus := []model.Sku{
		1625903,
		1148162,
		1076963,
		32638658,
		32605854,
		32205848,
		wrongSku,
	}
	transport := http.DefaultTransport

	httpClient := http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	productService := NewProductService(
		httpClient,
		"testToken",
		fmt.Sprintf("%s:%s", "http://localhost", "8082"),
		10,
	)

	products, err := productService.GetProductsBySkus(ctx, skus)
	require.ErrorIs(t, ErrProductNotFound, err)
	require.Equal(t, []model.Product(nil), products)
}
