package statistics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const metricNamespace = "robot"

func newPrometheusCounter(subSys, name, help string) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: subSys,
		Name:      name,
		Help:      help,
	})
}

func newPrometheusGauge(subSys, name, help string) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: metricNamespace,
		Subsystem: subSys,
		Name:      name,
		Help:      help,
	})
}

func newPrometheusSummary(subSys, name, help string, quantiles ...float64) prometheus.Summary {
	var objectives map[float64]float64
	if len(quantiles) > 0 {
		objectives = make(map[float64]float64, len(quantiles))
		for _, v := range quantiles {
			objectives[v] = 0
		}
	}

	return promauto.NewSummary(prometheus.SummaryOpts{
		Namespace:  metricNamespace,
		Subsystem:  subSys,
		Name:       name,
		Help:       help,
		Objectives: objectives,
	})
}

func newPrometheusHistogram(subSys, name, help string, buckets ...float64) prometheus.Histogram {
	return promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: metricNamespace,
		Subsystem: subSys,
		Name:      name,
		Help:      help,
		Buckets:   buckets,
	})
}
