package middlewares

import (
	"net/http"
	"route256/cart/internal/infra/logger"
)

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
