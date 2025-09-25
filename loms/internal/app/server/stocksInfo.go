package server

import (
	"context"
	"errors"
	lomspb "route256/loms/internal/api"
	"route256/loms/internal/domain/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) StocksInfo(ctx context.Context, request *lomspb.StocksInfoRequest) (*lomspb.StocksInfoResponse, error) {
	count, err := s.lomsServise.StocksInfo(ctx, request.Sku)
	if err != nil {
		if errors.Is(err, model.ErrShortOfStock) {
			return nil, status.Error(codes.FailedPrecondition, "")
		}
		return nil, status.Error(codes.Internal, "")
	}
	return &lomspb.StocksInfoResponse{

		Count: count,
	}, status.Errorf(codes.OK, "")
}
