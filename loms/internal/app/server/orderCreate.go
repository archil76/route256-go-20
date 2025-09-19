package server

import (
	"context"
	"errors"
	lomspb "route256/loms/internal/api"
	"route256/loms/internal/domain/model"
	lomsServise "route256/loms/internal/domain/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) OrderCreate(ctx context.Context, request *lomspb.OrderCreateRequest) (*lomspb.OrderCreateResponse, error) {
	items := []model.Item{}

	for _, reqItem := range request.Items {
		items = append(items, model.Item{
			Sku:   reqItem.Sku,
			Count: reqItem.Count})
	}

	orderID, err := s.lomsServise.OrderCreate(ctx, request.UserID, items)
	if err != nil {
		if errors.Is(err, lomsServise.ErrShortOfStock) {

			return nil, status.Errorf(codes.FailedPrecondition, "")
		}
		return nil, status.Error(codes.Internal, "")
	}

	return &lomspb.OrderCreateResponse{
		OrderID: orderID,
	}, status.Error(codes.OK, "")
}
