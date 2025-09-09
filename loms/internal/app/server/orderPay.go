package server

import (
	"context"
	lomspb "route256/loms/internal/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) OrderPay(ctx context.Context, request *lomspb.OrderPayRequest) (*lomspb.OrderPayResponse, error) {
	err := s.lomsServise.OrderPay(ctx, request.OrderId)
	if err != nil {
		return nil, err
	}

	return &lomspb.OrderPayResponse{}, status.Errorf(codes.OK, "")
}
