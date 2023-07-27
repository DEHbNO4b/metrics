package data

import (
	"errors"
	"strconv"
	"sync"
)

type Gauge struct {
	Name string
	Val  float64
}

type Counter struct {
	Name string
	Val  int64
}

type MetStore struct {
	Gauges   map[string]float64
	Counters map[string]int64
	sync.RWMutex
}

func NewMetStore() *MetStore {
	g := make(map[string]float64)
	c := make(map[string]int64)
	ms := MetStore{Gauges: g, Counters: c}
	return &ms
}
func (ms *MetStore) SetMetric(metric Metrics) error {
	switch metric.MType {
	case "gauge":
		ms.Lock()
		ms.Gauges[metric.ID] = *metric.Value
		ms.Unlock()
	case "counter":
		ms.Lock()
		ms.Counters[metric.ID] = ms.Counters[metric.ID] + *metric.Delta
		ms.Unlock()
	default:
		return errors.New("wrong metric type")
	}
	return nil
}
func (ms *MetStore) GetMetrics() []string {
	m := make([]string, 0, 20)
	for name, val := range ms.Gauges {
		m = append(m, name+":"+strconv.FormatFloat(val, 'f', -1, 64))
	}
	for name, val := range ms.Counters {
		m = append(m, name+":"+strconv.FormatInt(val, 10))
	}
	return m
}
func (ms *MetStore) SetGauge(g Gauge) error {
	// lock := sync.Mutex{}
	ms.Lock()
	ms.Gauges[g.Name] = g.Val
	ms.Unlock()
	return nil
}
func (ms *MetStore) SetCounter(c Counter) error {
	ms.Lock()
	ms.Counters[c.Name] = ms.Counters[c.Name] + c.Val
	ms.Unlock()
	return nil
}
func (ms *MetStore) GetGauge(name string) (Gauge, error) {
	g := Gauge{}
	g.Name = name
	ms.RLock()
	v, ok := ms.Gauges[name]
	ms.RUnlock()
	if !ok {
		return g, errors.New("not contains this metric")
	}
	g.Val = v
	return g, nil
}
func (ms *MetStore) GetCounter(name string) (Counter, error) {
	c := Counter{}
	c.Name = name
	ms.RLock()
	v, ok := ms.Counters[name]
	ms.RUnlock()
	if !ok {
		return c, errors.New("not contains this metric")
	}
	c.Val = v
	return c, nil
}
