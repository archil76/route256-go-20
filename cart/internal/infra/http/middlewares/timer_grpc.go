package middlewares

import (
	"context"
	"route256/cart/internal/infra/metrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func TimerUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	now := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)
	statusCode := int(status.Code(err)) //nolint:gosec
	duration := time.Since(now)
	metrics.StoreRequestDuration("", method, statusCode, duration)

	return err
}
