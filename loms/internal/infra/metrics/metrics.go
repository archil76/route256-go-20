package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "loms",
		Name:      "handler_request_total_counter",
		Help:      "Total count of request",
	}, []string{"protocol", "method", "path", "status"})

	requestDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "loms",
		Name:      "handler_request_duration_histogram",
		Help:      "Total duration of handler processing",
		Buckets:   prometheus.DefBuckets,
	}, []string{"protocol", "method", "path", "status"})
)

func IncRequestCount(protocol, method, path string, status int) {
	requestCounter.WithLabelValues(protocol, method, path, strconv.Itoa(status)).Inc()
}

func StoreRequestDuration(protocol, method, path string, status int, duration time.Duration) {
	requestDurationHistogram.WithLabelValues(protocol, method, path, strconv.Itoa(status)).Observe(float64(duration.Seconds()))
}
