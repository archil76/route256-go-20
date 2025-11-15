package middlewares

import (
	"net/http"
	"route256/cart/internal/infra/metrics"
	"time"
)

type CounterMux struct {
	h http.Handler
}

func NewCounterMux(h http.Handler) http.Handler {
	return &CounterMux{h: h}
}

func (m *CounterMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(now time.Time) {
		metrics.IncRequestCount(r.Method)
	}(time.Now())

	m.h.ServeHTTP(w, r)
}
