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
	SqlDB    interfaces.Database
	FileDB   interfaces.Database
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
	// m := data.Metrics{}
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

func (rs *MemStorage) GeMetricsData() []data.Metrics {
	metrics := make([]data.Metrics, 0, 30)
	for name, val := range rs.Gauges {
		m := data.NewMetric()
		m.ID = name
		m.MType = "gauge"
		*m.Value = val
		metrics = append(metrics, m)
	}
	for name, val := range rs.Counters {
		m := data.NewMetric()
		m.ID = name
		m.MType = "counter"
		*m.Delta = val
		metrics = append(metrics, m)
	}
	return metrics
}
func (rs *MemStorage) SetMetrics(m []data.Metrics) error {
	for _, el := range m {
		switch el.MType {
		case "gauge":
			rs.Gauges[el.ID] = *el.Value
		case "counter":
			rs.Counters[el.ID] = *el.Delta
		}
	}
	return nil
}
func (rs *MemStorage) StoreData() error {
	var db interfaces.Database
	if rs.SqlDB.Ping() {
		db = rs.SqlDB
	} else {
		db = rs.FileDB
	}
	data := rs.GeMetricsData()
	err := db.WriteMetrics(data)
	if err != nil {
		return err
	}
	return nil
}
