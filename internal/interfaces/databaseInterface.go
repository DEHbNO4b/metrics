package interfaces

import "github.com/DEHbNO4b/metrics/internal/data"

type Database interface {
	WriteMetrics([]data.Metrics) error
	ReadMetrics() ([]data.Metrics, error)
}
