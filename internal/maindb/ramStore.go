package maindb

import (
	"strconv"
	"sync"
	"time"

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
type StoreConfig struct {
	Filepath      string
	StoreInterval time.Duration
	Restore       bool
}
type RAMStore struct {
	Config   StoreConfig
	Gauges   map[string]float64
	Counters map[string]int64
	DB       interfaces.Database
	sync.RWMutex
}

func NewRAMStore(config StoreConfig) *RAMStore {
	g := make(map[string]float64)
	c := make(map[string]int64)
	rs := RAMStore{Config: config, Gauges: g, Counters: c}
	return &rs
}

func (rs *RAMStore) SetMetric(metric data.Metrics) error {
	switch metric.MType {
	case "gauge":
		rs.Lock()
		rs.Gauges[metric.ID] = *metric.Value
		rs.Unlock()
	case "counter":
		rs.Lock()
		rs.Counters[metric.ID] = rs.Counters[metric.ID] + *metric.Delta
		rs.Unlock()
	default:
		return interfaces.ErrWrongType
	}
	return nil
}
func (rs *RAMStore) GetMetrics() []string {
	m := make([]string, 0, 40)
	for name, val := range rs.Gauges {
		m = append(m, name+":"+strconv.FormatFloat(val, 'f', -1, 64))
	}
	for name, val := range rs.Counters {
		m = append(m, name+":"+strconv.FormatInt(val, 10))
	}
	return m
}
func (rs *RAMStore) GetMetric(met data.Metrics) (data.Metrics, error) {
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

func (rs *RAMStore) GeMetricsData() []data.Metrics {
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
func (rs *RAMStore) LoadFromStoreFile() error {

	metrics, err := rs.DB.ReadMetrics()
	if err != nil {
		// logger.Log.Error("unable to load data from DB", zap.Error(err))
		return err
	}
	for _, metric := range metrics {
		switch metric.MType {
		case "gauge":
			rs.Gauges[metric.ID] = *metric.Value
		case "counter":
			rs.Counters[metric.ID] = *metric.Delta
		}
	}
	return nil
}

// func (rs *RAMStore) StoreSchedule() {
// 	for {
// 		rs.StoreData()
// 		time.Sleep(rs.Config.StoreInterval)
// 	}
// }

func (rs *RAMStore) StoreData() error {
	data := rs.GeMetricsData()
	err := rs.DB.WriteMetrics(data)
	if err != nil {
		return err
	}
	return nil
}
