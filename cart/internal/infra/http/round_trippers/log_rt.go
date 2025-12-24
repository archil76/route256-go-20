package round_trippers

import (
	"net/http"
	"route256/cart/internal/infra/logger"
)

type LogRoundTripper struct {
	rt http.RoundTripper
}

func NewLogRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &LogRoundTripper{rt: rt}
}

func (l *LogRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	logger.Infow("service called",
		"url", r.URL.String(),
	)

	resp, err := l.rt.RoundTrip(r)
	if err != nil {
		logger.Infow("service called",
			"response", resp,
			"error", err,
		)
	}

	return resp, err
}
