package httpgateway

import (
	"strconv"

	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/oldjon/gx/service"
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	com.BasicMetrics
	loginTotal           *prometheus.CounterVec
	readHTTPReqFailTotal *prometheus.CounterVec
}

func newGatewayMetrics(driver service.ModuleDriver) *metrics {
	m := &metrics{}
	m.Metrics = driver.Metrics()
	m.MetricsNamespace = driver.HostName()
	m.MetricsSubSystem = driver.ModuleName()

	m.loginTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: m.MetricsNamespace,
			Subsystem: m.MetricsSubSystem,
			Name:      "platform_login_total",
			Help:      "total logins partitioned by platform type and server_id",
		},
		[]string{
			"platform",
			"login_type",
			"server_id",
		},
	)

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

	m.Metrics.MustRegister(m.loginTotal)
	m.Metrics.MustRegister(m.readHTTPReqFailTotal)
	return m
}

func (m *metrics) incLoginTotalMetrics(platform string, externalType, serverID uint32) { // nolint:unused
	m.loginTotal.With(prometheus.Labels{
		"platform":   platform,
		"login_type": strconv.Itoa(int(externalType)),
		"server_id":  strconv.Itoa(int(serverID)),
	}).Inc()
}

func (m *metrics) incReadHTTPFail(protocol string, err error) {
	m.readHTTPReqFailTotal.WithLabelValues(
		protocol,
		err.Error(),
	).Inc()
}
