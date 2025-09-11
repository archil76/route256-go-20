package server

import (
	"context"
	lomspb "route256/loms/internal/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) OrderCancel(ctx context.Context, request *lomspb.OrderCancelRequest) (*lomspb.OrderCancelResponse, error) {
	err := s.lomsServise.OrderCancel(ctx, request.OrderID)
	if err != nil {
		return nil, err
	}
	return &lomspb.OrderCancelResponse{}, status.Errorf(codes.OK, "")

}
