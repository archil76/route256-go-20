package middlewares

import (
	"context"
	"route256/loms/internal/infra/metrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func TimerUnaryServerInterceptor(ctx context.Context, request any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response any, err error) {
	now := time.Now()

	response, err = handler(ctx, request)

	duration := time.Since(now)

	statusCode := int(status.Code(err)) //nolint:gosec
	metrics.IncRequestCount("", info.FullMethod, statusCode)
	metrics.StoreRequestDuration("", info.FullMethod, statusCode, duration)

	return response, err
}
