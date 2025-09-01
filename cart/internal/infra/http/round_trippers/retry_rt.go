package round_trippers

import (
	"log"
	"net/http"
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
	log.Print("using product_service api")

	var response *http.Response
	var err error

	for attempts := 0; attempts < customRT.maxRetries; attempts++ {
		response, err = customRT.rt.RoundTrip(request)

		// good outcome
		if err == nil && !(response.StatusCode == http.StatusTooManyRequests || response.StatusCode == 420) {
			break
		}

		// delay and retry
		select {
		case <-request.Context().Done():
			return response, request.Context().Err()
		case <-time.After(customRT.delay):
		}
	}

	return response, err
}
