package interfaces

import (
	"errors"

	"github.com/DEHbNO4b/metrics/internal/data"
)

var ErrNotContains error = errors.New("not contains this metric")
var ErrWrongType error = errors.New("wrong metric type")

// MetricsStorage it is an interface for storing metrics in RAM.
type MetricsStorage interface {
	GetMetrics() []data.Metrics
	SetMetric(data.Metrics) error
	GetMetric(data.Metrics) (data.Metrics, error)
}
