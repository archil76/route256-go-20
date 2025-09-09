package server

import (
	"context"
	lomspb "route256/loms/internal/api"
	"route256/loms/internal/domain/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) OrderCreate(ctx context.Context, request *lomspb.OrderCreateRequest) (*lomspb.OrderCreateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Невалидный запрос")
	}

	items := []model.Item{}

	for _, reqItem := range request.Items {
		items = append(items, model.Item{
			Sku:   reqItem.Sku,
			Count: reqItem.Count})
	}

	orderId, err := s.lomsServise.OrderCreate(ctx, request.UserId, items)
	if err != nil {
		return nil, err
	}

	return &lomspb.OrderCreateResponse{
		OrderId: orderId,
	}, status.Errorf(codes.OK, "")
}
