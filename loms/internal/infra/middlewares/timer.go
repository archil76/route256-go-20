package middlewares

import (
	"net/http"
	"route256/loms/internal/infra/logger"
	"route256/loms/internal/infra/metrics"
	"time"

	"go.uber.org/zap"
)

type TimerMux struct {
	h http.Handler
}

func NewTimerMux(h http.Handler) http.Handler {
	return &TimerMux{h: h}
}

func (m *TimerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	rw := newResponseWriter(w)
	m.h.ServeHTTP(rw, r)

	duration := time.Since(now)

	metrics.StoreRequestDuration("http", r.Method, r.Pattern, rw.statusCode, duration)
	logger.Infow("handler spent time", r.Method, r.URL.Path, zap.String("mc", duration.String()))
}
