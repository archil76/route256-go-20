package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	queriesCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "loms",
		Name:      "db_query_total_counter",
		Help:      "Total count of query",
	}, []string{"queryType"})

	queriesDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "loms",
		Name:      "db_query_duration_histogram",
		Help:      "Total query execution time",
		Buckets:   prometheus.DefBuckets,
	}, []string{"queryType", "error"})
)

func IncQueryCount(queryType string) {
	queriesCounter.WithLabelValues(queryType).Inc()
}

func StoreQueryDuration(queryType, errorCode string, duration time.Duration) {
	queriesDurationHistogram.WithLabelValues(queryType, errorCode).Observe(float64(duration.Seconds()))
}
