package server

import (
	"context"
	lomspb "route256/loms/internal/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) StocksInfo(ctx context.Context, request *lomspb.StocksInfoRequest) (*lomspb.StocksInfoResponse, error) {
	count, err := s.lomsServise.StocksInfo(ctx, request.Sku)
	if err != nil {
		return nil, err
	}
	return &lomspb.StocksInfoResponse{

		Count: count,
	}, status.Errorf(codes.OK, "")
}
