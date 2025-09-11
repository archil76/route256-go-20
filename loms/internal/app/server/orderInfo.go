package server

import (
	"context"
	desc "route256/loms/internal/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) OrderInfo(ctx context.Context, request *desc.OrderInfoRequest) (*desc.OrderInfoResponse, error) {
	order, err := s.lomsServise.OrderInfo(ctx, request.OrderId)
	if err != nil {
		return nil, err
	}

	orderInfoResponse := desc.OrderInfoResponse{
		Status:  string(order.Status),
		OrderId: order.OrderID,
		Items:   []*desc.Items{},
	}
	orderInfoResponse.OrderId = order.OrderID
	for _, item := range order.Items {

		orderInfoResponse.Items = append(orderInfoResponse.Items, &desc.Items{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	return &orderInfoResponse, status.Errorf(codes.OK, "")
}
