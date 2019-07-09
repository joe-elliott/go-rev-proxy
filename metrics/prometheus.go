package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "cert_exporter"
)

var (
	RequestLatencyMilliseconds = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "requests_latency_ms",
			Help:       "Summary of Request Latency (ms)",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(RequestLatencyMilliseconds)
}
