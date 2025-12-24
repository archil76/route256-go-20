package round_trippers

import (
	"net/http"
	"route256/cart/internal/infra/metrics"
)

type CounterRoundTripper struct {
	rt http.RoundTripper
}

func NewCounterRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &CounterRoundTripper{rt: rt}
}

func (l *CounterRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	statusCode := 404

	resp, err := l.rt.RoundTrip(r)
	if resp != nil {
		statusCode = resp.StatusCode
	}
	metrics.IncRequestCount(r.Method, r.Pattern, statusCode)

	return resp, err
}
