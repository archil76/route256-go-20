package round_trippers

import (
	"net/http"
	"route256/cart/internal/infra/metrics"
	"time"
)

type TimerRoundTripper struct {
	rt http.RoundTripper
}

func NewTimerRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &CounterRoundTripper{rt: rt}
}

func (l *TimerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	statusCode := 404
	now := time.Now()

	resp, err := l.rt.RoundTrip(r)
	if resp != nil {
		statusCode = resp.StatusCode
	}

	duration := time.Since(now)

	metrics.StoreRequestDuration(r.Method, r.Pattern, statusCode, duration)

	return resp, err
}
