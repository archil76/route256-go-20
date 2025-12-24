package middlewares

import (
	"context"
	"route256/cart/internal/infra/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func CounterUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	err := invoker(ctx, method, req, reply, cc, opts...)
	statusCode := int(status.Code(err)) //nolint:gosec

	metrics.IncRequestCount("", method, statusCode)
	return err
}
