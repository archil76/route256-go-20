package round_trippers

import (
	"net/http"
	"route256/cart/internal/infra/logger"
	"time"
)

type RetryRoundTripper struct {
	rt         http.RoundTripper
	maxRetries int
	delay      time.Duration
}

func NewRetryRoundTripper(rt http.RoundTripper, maxRetries int, delay time.Duration) http.RoundTripper {
	return &RetryRoundTripper{rt: rt, maxRetries: maxRetries, delay: delay}
}

func (customRT *RetryRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error
	var attempts int
	var resStatusCode int

	for attempts = 0; attempts < customRT.maxRetries; attempts++ {
		response, err = customRT.rt.RoundTrip(request)
		resStatusCode = response.StatusCode

		// good outcome
		if err == nil && !(resStatusCode == http.StatusTooManyRequests || resStatusCode == 420) {
			break
		}

		// delay and retry
		select {
		case <-request.Context().Done():
			return response, request.Context().Err()
		case <-time.After(customRT.delay):
		}
	}

	logger.Infow("using product_service api",
		"url", request.URL.String(),
		"method", request.Method,
		"attempts", attempts,
		"resStatusCode", resStatusCode,
		"error", err)

	return response, err
}
