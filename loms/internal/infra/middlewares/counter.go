package middlewares

import (
	"net/http"
	"route256/loms/internal/infra/metrics"
)

type CounterMux struct {
	h http.Handler
}

func NewCounterMux(h http.Handler) http.Handler {
	return &CounterMux{h: h}
}

func (m *CounterMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := newResponseWriter(w)

	m.h.ServeHTTP(rw, r)

	metrics.IncRequestCount("http", r.Method, r.Pattern, rw.statusCode)
}
