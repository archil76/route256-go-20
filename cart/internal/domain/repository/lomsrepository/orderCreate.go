package lomsrepository

import (
	"context"
	desc "route256/cart/internal/api"

	"route256/cart/internal/domain/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *LomsService) OrderCreate(ctx context.Context, userID model.UserID, reportCart *model.ReportCart) (int64, error) {
	conn, err := grpc.NewClient(
		s.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return 0, err
	}

	client := desc.NewLomsClient(conn)

	orderCreateRequest := desc.OrderCreateRequest{UserID: userID}
	for _, itemInCart := range reportCart.Items {
		orderCreateRequest.Items = append(orderCreateRequest.Items, &desc.Items{
			Sku:   itemInCart.SKU,
			Count: itemInCart.Count,
		})
	}

	orderCreateResponse, err := client.OrderCreate(ctx, &orderCreateRequest)
	if err != nil {
		return 0, ErrOrderNotFound
	}

	return orderCreateResponse.OrderID, nil
}
