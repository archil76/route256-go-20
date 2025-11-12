package middlewares

import (
	"net/http"
	"route256/cart/internal/infra/logger"
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
	now := time.Now()
	m.h.ServeHTTP(w, r)
	logger.Infow("handler spent time", zap.String("mc", time.Since(now).String()))
}
