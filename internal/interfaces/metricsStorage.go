package interfaces

import "github.com/DEHbNO4b/metrics/internal/data"

type MetricsStorage interface {
	GetMetrics() []string
	SetMetric(data.Metrics) error
	GetMetric(data.Metrics) (data.Metrics, error)
}
