package interfaces

import "github.com/DEHbNO4b/metrics/internal/data"

type MetricsStorage interface {
	GetMetrics() []string
	SetMetric(data.Metrics) error
	// SetGauge(g data.Gauge) error
	// SetCounter(c data.Counter) error
	GetGauge(name string) (data.Gauge, error)
	GetCounter(name string) (data.Counter, error)

	GetMetric(data.Metrics) (data.Metrics, error)
}
