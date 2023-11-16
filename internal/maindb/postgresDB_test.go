package maindb

import (
	"testing"

	"github.com/DEHbNO4b/metrics/internal/data"
)

func BenchmarkPostgresWriteMetrics(b *testing.B) {
	p, _ := NewPostgresDB("postgres://postgres:917836@localhost:5432/test?sslmode=disable")
	metricsData := make([]data.Metrics, 0, 3)
	metr := data.NewMetric()
	metr.ID = "test"
	metr.MType = "gauge"
	*metr.Value = 3.14
	metricsData = append(metricsData, metr)

	for i := 0; i < b.N; i++ {
		p.WriteMetrics(metricsData)
	}
}
func BenchmarkPostgresReadMetrics(b *testing.B) {
	p, _ := NewPostgresDB("postgres://postgres:917836@localhost:5432/test?sslmode=disable")
	// for _, size := range []int{1, 10, 100, 1000, 10000} {
	metricsData := getRundomMetrics(20)
	p.WriteMetrics(metricsData)

	b.ResetTimer()

	// name := fmt.Sprintf("Contains-%d", size)
	// b.Run(name, func(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p.ReadMetrics()
	}
	// })
	// }
}
