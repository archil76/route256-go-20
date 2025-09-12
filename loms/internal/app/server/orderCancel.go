package server

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	lomspb "route256/loms/internal/api"
	lomsServise "route256/loms/internal/domain/service"
)

func (s Server) OrderCancel(ctx context.Context, request *lomspb.OrderCancelRequest) (*lomspb.OrderCancelResponse, error) {
	err := s.lomsServise.OrderCancel(ctx, request.OrderID)
	if err != nil {
		if errors.Is(err, lomsServise.ErrOrderDoesntExist) {
			return nil, status.Error(codes.NotFound, "")
		}

		if errors.Is(err, lomsServise.ErrInvalidOrderStatus) {
			return nil, status.Error(codes.FailedPrecondition, "")
		}
		return nil, status.Error(codes.Internal, "")
	}
	return &lomspb.OrderCancelResponse{}, status.Errorf(codes.OK, "")

}
