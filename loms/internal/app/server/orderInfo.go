package server

import (
	"context"
	"errors"
	desc "route256/loms/internal/api"
	lomsServise "route256/loms/internal/domain/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) OrderInfo(ctx context.Context, request *desc.OrderInfoRequest) (*desc.OrderInfoResponse, error) {
	order, err := s.lomsServise.OrderInfo(ctx, request.OrderID)
	if err != nil {
		if errors.Is(err, lomsServise.ErrOrderDoesntExist) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	orderInfoResponse := desc.OrderInfoResponse{
		Status:  string(order.Status),
		OrderID: order.OrderID,
		Items:   []*desc.Items{},
	}
	orderInfoResponse.OrderID = order.OrderID
	for _, item := range order.Items {

		orderInfoResponse.Items = append(orderInfoResponse.Items, &desc.Items{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	return &orderInfoResponse, status.Errorf(codes.OK, "")
}
