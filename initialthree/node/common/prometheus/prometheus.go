package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

/*
- Counter    一种累加的 metric，典型的应用如：请求的个数，结束的任务数， 出现的错误数等等.
- Gauge      一种常规的 metric，典型的应用如：温度，运行的 goroutines 的个数。可以任意加减。
- Histogram  可以理解为柱状图，典型的应用如：请求持续时间，响应大小。可以对观察结果采样，分组及统计。
- Summary    类似于 Histogram, 典型的应用如：请求持续时间，响应大小。提供观测值的 count 和 sum 功能。提供百分位的功能，即可以按百分比划分跟踪结果。
*/

func NewPrometheusCounter(namespace, subSys, name, help string) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subSys,
		Name:      name,
		Help:      help,
	})
}

func NewPrometheusGauge(namespace, subSys, name, help string) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subSys,
		Name:      name,
		Help:      help,
	})
}

func NewPrometheusGaugeFunc(namespace, subSys, name, help string, f func() float64) prometheus.GaugeFunc {
	return promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subSys,
		Name:      name,
		Help:      help,
	}, f)
}

func NewPrometheusGaugeVec(namespace, subSys, name, help string, labelNames []string) *prometheus.GaugeVec {
	return promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subSys,
		Name:      name,
		Help:      help,
	}, labelNames)
}

func NewPrometheusSummary(namespace, subSys, name, help string, quantiles ...float64) prometheus.Summary {
	var objectives map[float64]float64
	if len(quantiles) > 0 {
		objectives = make(map[float64]float64, len(quantiles))
		for _, v := range quantiles {
			objectives[v] = 0
		}
	}

	return promauto.NewSummary(prometheus.SummaryOpts{
		Namespace:  namespace,
		Subsystem:  subSys,
		Name:       name,
		Help:       help,
		Objectives: objectives,
	})
}

func NewPrometheusHistogram(namespace, subSys, name, help string, buckets ...float64) prometheus.Histogram {
	return promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subSys,
		Name:      name,
		Help:      help,
		Buckets:   buckets,
	})
}
