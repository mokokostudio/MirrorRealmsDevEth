package gmgateway

import (
	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/oldjon/gx/service"
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	com.BasicMetrics
	readHTTPReqFailTotal *prometheus.CounterVec
}

func newGatewayMetrics(driver service.ModuleDriver) *metrics {
	m := &metrics{}
	m.Metrics = driver.Metrics()
	m.MetricsNamespace = driver.HostName()
	m.MetricsSubSystem = driver.ModuleName()

	m.readHTTPReqFailTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.MetricsNamespace,
			Subsystem: m.MetricsSubSystem,
			Name:      "read_http_req_fail_total",
			Help:      "read http req fail total",
		},
		[]string{
			"protocol",
			"error",
		},
	)

	m.Metrics.MustRegister(m.readHTTPReqFailTotal)
	return m
}

func (m *metrics) incReadHTTPFail(protocol string, err error) {
	m.readHTTPReqFailTotal.WithLabelValues(
		protocol,
		err.Error(),
	).Inc()
}
