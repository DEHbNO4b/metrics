package interfaces

import "github.com/DEHbNO4b/metrics/internal/data"

type SourceInterface interface {
	GeMetricsData() []data.Metrics
}
