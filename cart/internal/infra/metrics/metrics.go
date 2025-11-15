package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "cart",
		Name:      "handler_request_total_counter",
		Help:      "Total count of request",
	}, []string{"handler"})

	requestDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cart",
		Name:      "handler_request_duration_histogram",
		Help:      "Total duration of handler processing",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler"})

	repoSizeGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "cart",
		Name:      "repo_size_gauge",
		Help:      "Size of repo",
	})
)

func IncRequestCount(handler string) {
	requestCounter.WithLabelValues(handler).Inc()
}

func StoreRequestDuration(handler string, duration time.Duration) {
	requestDurationHistogram.WithLabelValues(handler).Observe(float64(duration.Seconds()))
}

func StoreRepoSize(size float64) {
	repoSizeGauge.Set(size)
}
