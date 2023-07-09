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

func (g Gauge) String() string {
	return g.Name + ": " + strconv.FormatFloat(g.Val, 'f', -1, 64)
}

type Counter struct {
	Name string
	Val  int64
}

func (c Counter) String() string {
	return c.Name + ": " + strconv.FormatInt(c.Val, 10)
}

type MetStore struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMetStore() *MetStore {
	g := make(map[string]float64)
	c := make(map[string]int64)
	ms := MetStore{Gauges: g, Counters: c}
	return &ms
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
	lock := sync.Mutex{}
	lock.Lock()
	ms.Gauges[g.Name] = g.Val
	lock.Unlock()
	return nil
}
func (ms *MetStore) SetCounter(c Counter) error {

	ms.Counters[c.Name] = ms.Counters[c.Name] + c.Val
	return nil
}
func (ms *MetStore) GetGauge(name string) (Gauge, error) {
	g := Gauge{}
	g.Name = name
	v, ok := ms.Gauges[name]
	if !ok {
		return g, errors.New("not contains this metric")
	}
	g.Val = v
	return g, nil
}
func (ms *MetStore) GetCounter(name string) (Counter, error) {
	c := Counter{}
	c.Name = name
	v, ok := ms.Counters[name]
	if !ok {
		return c, errors.New("not contains this metric")
	}
	c.Val = v
	return c, nil
}
