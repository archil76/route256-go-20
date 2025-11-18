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
	}, []string{"method", "path", "status"})

	requestDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "loms",
		Name:      "handler_request_duration_histogram",
		Help:      "Total duration of handler processing",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	repoSizeGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "loms",
		Name:      "repo_size_gauge",
		Help:      "Size of repo",
	})
)

func IncRequestCount(method, path string, status int) {
	requestCounter.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
}

func StoreRequestDuration(method, path string, status int, duration time.Duration) {
	requestDurationHistogram.WithLabelValues(method, path, strconv.Itoa(status)).Observe(float64(duration.Seconds()))
}

func StoreRepoSize(size float64) {
	repoSizeGauge.Set(size)
}
