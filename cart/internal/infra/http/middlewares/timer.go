package middlewares

import (
	"net/http"
	"route256/cart/internal/infra/logger"
	"route256/cart/internal/infra/metrics"
	"time"

	"go.uber.org/zap"
)

type TimerMux struct {
	h http.Handler
}

func NewTimeMux(h http.Handler) http.Handler {
	return &TimerMux{h: h}
}

func (m *TimerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/metrics" {
		m.h.ServeHTTP(w, r)
		return
	}
	now := time.Now()

	rw := newResponseWriter(w)
	m.h.ServeHTTP(rw, r)

	duration := time.Since(now)

	metrics.StoreRequestDuration(r.Method, r.Pattern, rw.statusCode, duration)
	logger.Infow("handler spent time", r.Method, r.URL.Path, zap.String("mc", duration.String()))
}
