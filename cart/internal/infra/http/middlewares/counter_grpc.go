package middlewares

import (
	"context"
	"route256/cart/internal/infra/metrics"

	"google.golang.org/grpc"
)

func CounterUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	err := invoker(ctx, method, req, reply, cc, opts...)
	metrics.IncRequestCount(method, "req.Pattern", 444)
	return err
}
