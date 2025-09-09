package server

import (
	"context"
	lomspb "route256/loms/internal/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) OrderCancel(ctx context.Context, request *lomspb.OrderCancelRequest) (*lomspb.OrderCancelResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Невалидный запрос")
	}

	err := s.lomsServise.OrderCancel(ctx, request.OrderId)
	if err != nil {
		return nil, err
	}
	return &lomspb.OrderCancelResponse{}, status.Errorf(codes.OK, "")

}
