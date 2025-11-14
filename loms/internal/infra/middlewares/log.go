package middlewares

import (
	"context"
	"net/http"
	"route256/loms/internal/infra/logger"

	"google.golang.org/grpc"
)

func Log(ctx context.Context, request any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response any, err error) {
	logger.Infow("request got", "Method", info.FullMethod)

	response, err = handler(ctx, request)
	if err != nil {
		logger.Errorw("error in handler", "Method", info.FullMethod, "Error", err)
	}
	logger.Infow("request processed", "Method", info.FullMethod)

	return response, err
}

type LogMux struct {
	h http.Handler
}

func NewLogMux(h http.Handler) http.Handler {
	return &LogMux{h: h}
}

func (m *LogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Infow("request got", r.Method, r.URL.Path)

	m.h.ServeHTTP(w, r)

	logger.Infow("request processed", r.Method, r.URL.Path)
}
