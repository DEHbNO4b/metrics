package maindb

import (
	"math/rand"

	"github.com/DEHbNO4b/metrics/internal/data"
)

func getRundomMetrics(num int) []data.Metrics {
	letters := []byte("abcdefghigklmnopqrstuvwxyz")
	src := rand.New(rand.NewSource(25))
	metrics := make([]data.Metrics, 0, 30)
	for i := 0; i < num; i++ {
		m := data.NewMetric()
		var str []byte
		for k := 0; k < 6; k++ {
			str = append(str, letters[src.Intn(26)])
		}
		m.ID = string(str)
		m.MType = "gauge"
		*m.Value = src.Float64()
		metrics = append(metrics, m)
	}
	return metrics
}
