package maindb

import (
	"sync"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/interfaces"
)

type Gauge struct {
	Name string
	Val  float64
}

type Counter struct {
	Name string
	Val  int64
}

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
	sync.RWMutex
}

func NewMemStorage() *MemStorage {
	g := make(map[string]float64)
	c := make(map[string]int64)
	rs := MemStorage{Gauges: g, Counters: c}
	return &rs
}

func (rs *MemStorage) SetMetric(metric data.Metrics) error {
	switch metric.MType {
	case "gauge":
		rs.Lock()
		rs.Gauges[metric.ID] = *metric.Value
		rs.Unlock()
	case "counter":
		rs.Lock()
		rs.Counters[metric.ID] = *metric.Delta
		rs.Unlock()
	default:
		return interfaces.ErrWrongType
	}
	return nil
}
func (rs *MemStorage) GetMetrics() []data.Metrics {
	metrics := make([]data.Metrics, 0, 30)
	for name, val := range rs.Gauges {
		metric := data.NewMetric()
		metric.ID = name
		metric.MType = "gauge"
		*metric.Value = val
		metrics = append(metrics, metric)
	}
	for name, val := range rs.Counters {
		metric := data.NewMetric()
		metric.ID = name
		metric.MType = "counter"
		*metric.Delta = val
		metrics = append(metrics, metric)

	}
	return metrics
}
func (rs *MemStorage) GetMetric(met data.Metrics) (data.Metrics, error) {
	switch met.MType {
	case "gauge":
		val, ok := rs.Gauges[met.ID]
		if !ok {
			return data.Metrics{}, interfaces.ErrNotContains
		}
		met.Value = &val
	case "counter":
		del, ok := rs.Counters[met.ID]
		if !ok {
			return data.Metrics{}, interfaces.ErrNotContains
		}
		met.Delta = &del
	default:
		return data.Metrics{}, interfaces.ErrWrongType
	}

	return met, nil
}
