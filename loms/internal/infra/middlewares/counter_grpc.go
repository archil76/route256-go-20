package middlewares

import (
	"context"
	"route256/loms/internal/infra/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func CounterUnaryServerInterceptor(ctx context.Context, request any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response any, err error) {
	response, err = handler(ctx, request)

	statusCode := int(status.Code(err)) //nolint:gosec
	metrics.IncRequestCount("grpc", "", info.FullMethod, statusCode)

	return response, err
}
