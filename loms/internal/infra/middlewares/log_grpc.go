package middlewares

import (
	"context"
	"route256/loms/internal/infra/logger"

	"google.golang.org/grpc"
)

func LogUnaryServerInterceptor(ctx context.Context, request any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response any, err error) {
	logger.Infow("request got", "Method", info.FullMethod)

	response, err = handler(ctx, request)
	if err != nil {
		logger.Errorw("error in handler", "Method", info.FullMethod, "Error", err)
	}
	logger.Infow("request processed", "Method", info.FullMethod)

	return response, err
}
