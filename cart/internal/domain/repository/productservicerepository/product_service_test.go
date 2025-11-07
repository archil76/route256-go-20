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

type testDataStruct struct {
	name      string
	userID    int64
	item      model.Item
	wantedErr error
}

func TestHandler_All(t *testing.T) {
	t.Run("TestHandler_AddItem_Table", TestHandler_Product_Service)

}

func TestHandler_Product_Service(t *testing.T) {
	ctx := context.Background()
	skus := []model.Sku{
		1625903,
		1148162,
		1076963,
		32638658,
		32605854,
		32205848,
		32205849,
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
	)

	products, err := productService.GetProductsBySkus(ctx, skus)
	require.NoError(t, err)
	require.Equal(t, len(skus), len(products))
}
