package util

import (
	"errors"
	"strconv"

	"github.com/oldjon/gx/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type ServiceMetricsI interface {
	UpdateResLoadFileNumMetrics()
	UpdateResLoadFailedTimesMetrics(err *ResMetricError)
}

type ServiceMetrics struct {
	Logger           *zap.Logger
	Metrics          prometheus.Registerer
	MetricsNamespace string
	MetricsSubSystem string

	ResLoadFailedTimes *prometheus.CounterVec
	ResLoadFileNum     *prometheus.CounterVec
}

func NewServiceMetrics(driver service.ModuleDriver) *ServiceMetrics {
	mt := &ServiceMetrics{}
	mt.InitCommonMetrics(driver)
	return mt
}

func (sm *ServiceMetrics) InitCommonMetrics(driver service.ModuleDriver) {
	sm.Logger = driver.Logger()
	host := driver.Host()
	sm.Metrics = host.Metrics()
	sm.MetricsNamespace = host.Name()
	sm.MetricsSubSystem = driver.ModuleName()
	sm.ResLoadFailedTimes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: sm.MetricsNamespace,
			Name:      "res_load_failed_times",
			Help:      "res load failed times",
		},
		[]string{
			"service",
			"csv",
			"field",
			"row",
			"error",
		},
	)
	sm.ResLoadFileNum = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: sm.MetricsNamespace,
			Name:      "res_load_file_num",
			Help:      "res load file num",
		},
		[]string{
			"service",
		},
	)
	sm.ResLoadFailedTimes = registerOrGet(sm.Metrics, sm.ResLoadFailedTimes).(*prometheus.CounterVec)
	sm.ResLoadFileNum = registerOrGet(sm.Metrics, sm.ResLoadFileNum).(*prometheus.CounterVec)
}

func (sm *ServiceMetrics) UpdateResLoadFailedTimesMetrics(err *ResMetricError) {
	sm.ResLoadFailedTimes.With(prometheus.Labels{
		"service": sm.MetricsSubSystem,
		"csv":     err.ResName,
		"field":   err.ResField,
		"row":     strconv.Itoa(err.ResRow),
		"error":   err.Error(),
	}).Inc()
}

func (sm *ServiceMetrics) UpdateResLoadFileNumMetrics() {
	sm.ResLoadFileNum.With(prometheus.Labels{
		"service": sm.MetricsSubSystem,
	}).Inc()
}

type ResMetricError struct {
	ResName  string
	ResField string
	ResRow   int
	Err      error
}

func NewResMetricError(csvName, field string, row int, err error) error {
	return &ResMetricError{
		ResName:  csvName,
		ResField: field,
		ResRow:   row + 1,
		Err:      err,
	}
}

func (r ResMetricError) Error() string {
	if r.Err == nil {
		return ""
	}
	return r.Err.Error()
}

func BuildResLoadFuncWithMetrics(sm ServiceMetricsI, fn func(string) error) func(string) error {
	return func(s string) error {
		sm.UpdateResLoadFileNumMetrics()
		err := fn(s)
		if err != nil {
			var rErr *ResMetricError
			if errors.As(err, &rErr) {
				sm.UpdateResLoadFailedTimesMetrics(rErr)
				err = rErr.Err
			} else {
				sm.UpdateResLoadFailedTimesMetrics(&ResMetricError{ResName: s, Err: err})
			}
		}
		return err
	}
}

func registerOrGet(r prometheus.Registerer, c prometheus.Collector) prometheus.Collector {
	if err := r.Register(c); err != nil {
		var nErr prometheus.AlreadyRegisteredError
		if errors.As(err, &nErr) {
			return nErr.ExistingCollector
		}
		panic(err)
	}
	return c
}
