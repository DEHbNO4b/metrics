// package interfaces contains interfaces.
package interfaces

import "github.com/DEHbNO4b/metrics/internal/data"

// Database it is an interface for storing metrics.
type Database interface {
	WriteMetrics([]data.Metrics) error
	ReadMetrics() ([]data.Metrics, error)
}
