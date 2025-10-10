package server

import (
	"context"
	"errors"
	lomspb "route256/loms/internal/api"
	"route256/loms/internal/domain/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) OrderPay(ctx context.Context, request *lomspb.OrderPayRequest) (*lomspb.OrderPayResponse, error) {
	err := s.lomsServise.OrderPay(ctx, request.OrderID)
	if err != nil {
		if errors.Is(err, model.ErrOrderDoesntExist) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else if errors.Is(err, model.ErrInvalidOrderStatus) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Error(codes.Unknown, "")
	}

	return &lomspb.OrderPayResponse{}, status.Errorf(codes.OK, "")
}
