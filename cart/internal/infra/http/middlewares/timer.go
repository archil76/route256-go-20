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
	defer func(now time.Time) {
		duration := time.Since(now)
		metrics.StoreRequestDuration(r.Method, duration)
		logger.Infow("handler spent time", zap.String("mc", duration.String()))
	}(time.Now())
	m.h.ServeHTTP(w, r)
}
