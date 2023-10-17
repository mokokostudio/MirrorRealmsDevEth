package common

import "github.com/prometheus/client_golang/prometheus"

type BasicMetrics struct {
	Metrics          prometheus.Registerer
	MetricsNamespace string
	MetricsSubSystem string
}
